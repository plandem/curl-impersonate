package curl

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/plandem/curl-impersonate/presets"
	"github.com/plandem/curl-impersonate/types"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Option func(o *Curl)

type Curl struct {
	headers *types.Headers
	flags   *types.Flags
	preset  presets.Preset
	binary  string
	isValid bool
}

func New(options ...Option) *Curl {
	curl := &Curl{
		preset:  presets.Default(),
		headers: types.NewHeaders(),
		flags: types.NewFlags(
			types.Flag("include", true),
			types.Flag("silent", true),
		),
	}

	curl.Set(options...)

	if curl.binary == "" {
		curl.binary = "curl"
	}

	return curl
}

func (curl *Curl) Set(options ...Option) {
	for _, o := range options {
		o(curl)
	}
}

func Header(key string, value string) func(*Curl) {
	return func(curl *Curl) {
		curl.headers.Set(key, value)
	}
}

func Flag(name string, value interface{}) func(*Curl) {
	return func(curl *Curl) {
		curl.flags.Set(name, value)
	}
}

func Preset(preset presets.PresetFn) func(*Curl) {
	return func(curl *Curl) {
		curl.preset = preset()
	}
}

func Binary(fileName string) func(*Curl) {
	return func(curl *Curl) {
		curl.binary = fileName
		curl.isValid = false
	}
}

func (curl *Curl) Validate() error {
	curl.isValid = false

	if _, err := os.Stat(curl.binary); os.IsNotExist(err) {
		if fullPath, err := exec.LookPath(curl.binary); err != nil {
			return err
		} else {
			curl.binary = fullPath
		}
	}

	curl.isValid = true
	return nil
}

func (curl *Curl) SetHeader(key string, value string) {
	curl.headers.Set(key, value)
}

func (curl *Curl) SetFlag(name string, value interface{}) {
	curl.flags.Set(name, value)
}

// Request return http.Response to make it similar to other fetching libraries, but also returns underlying []byte of content for simpler usage without extra allocation to read resp.Body
func (curl *Curl) Request(url string) (*http.Response, []http.Header, []byte, error) {
	if !curl.isValid {
		if err := curl.Validate(); err != nil {
			return nil, nil, nil, err
		}
	}

	var args []string

	args = append(args, curl.preset.Headers.Generate(false)...)
	args = append(args, curl.headers.Generate(true)...)
	args = append(args, curl.preset.Flags.Generate()...)
	args = append(args, curl.flags.Generate()...)
	args = append(args, url)

	cmd := exec.Command(curl.binary, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode := exitErr.ExitCode()
			return nil, nil, nil, &Error{ExitCode: exitCode}
		}

		return nil, nil, nil, fmt.Errorf("unexpected error executing curl: %w. stderr: %s", err, stderr.String())
	}

	// Get the full output from stdout
	output := stdout.Bytes()

	responses, lastBody, err := extractAllResponses(output)
	if err != nil {
		return nil, nil, nil, err
	}

	var headers []http.Header
	for _, rawHeaders := range responses {
		header := parseHeaders(rawHeaders)
		if len(header) > 0 {
			headers = append(headers, header)
		}
	}

	lastHeaders := responses[len(responses)-1]
	statusLine := bytes.Split(lastHeaders, []byte("\n"))[0]
	statusParts := bytes.Split(statusLine, []byte(" "))
	if len(statusParts) < 2 {
		return nil, nil, nil, errors.New("unable to extract status code")
	}

	statusCode, err := getStatusCode(string(statusParts[1]))
	if err != nil {
		return nil, nil, nil, err
	}
	if statusCode >= 400 {
		return nil, nil, nil, &HTTPError{
			StatusCode: statusCode,
			Status:     http.StatusText(statusCode),
		}
	}

	return &http.Response{
		StatusCode: statusCode,
		Header:     headers[len(headers)-1],
		Body:       io.NopCloser(bytes.NewReader(lastBody)),
	}, headers, lastBody, nil
}

// extract headers/body respecting multi-responses (e.g. proxy + redirect)
func extractAllResponses(output []byte) ([][]byte, []byte, error) {
	var headers [][]byte
	var lastBody []byte
	lastHeaderIndex := -1

	// Track the current position in the original buffer
	pos := 0

	for {
		// Find the next "\r\n\r\n" delimiter
		end := bytes.Index(output[pos:], []byte("\r\n\r\n"))
		if end == -1 {
			break
		}

		// Extract the current block
		block := output[pos : pos+end]

		// Check if the block starts with "HTTP/"
		if bytes.HasPrefix(block, []byte("HTTP/")) && len(bytes.TrimSpace(block)) > len("HTTP/") {
			headers = append(headers, block)
			lastHeaderIndex = len(headers) - 1
		}

		// Move the position past the current block and the delimiter
		pos += end + len("\r\n\r\n")
	}

	// If a last valid HTTP response block was found, capture the body content
	if lastHeaderIndex != -1 {
		lastBody = output[pos:]
	}

	if len(headers) == 0 {
		return nil, nil, errors.New("unable to extract valid HTTP responses")
	}

	return headers, lastBody, nil
}

func getStatusCode(statusCode string) (int, error) {
	code, err := strconv.Atoi(statusCode)
	if err != nil {
		return 0, fmt.Errorf("invalid status code format: %s", statusCode)
	}
	return code, nil
}

func parseHeaders(headers []byte) http.Header {
	result := make(http.Header)
	lines := bytes.Split(headers, []byte("\r\n"))
	for _, line := range lines {
		if bytes.Contains(line, []byte(":")) {
			parts := bytes.SplitN(line, []byte(":"), 2)
			key := strings.TrimSpace(string(parts[0]))
			value := strings.TrimSpace(string(parts[1]))
			result.Add(key, value)
		}
	}
	return result
}

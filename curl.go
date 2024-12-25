package curl

import (
	"binance-news/curl/presets"
	"binance-news/curl/types"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Option func(o *Curl)

type HTTPError struct {
	StatusCode int
	Status     string
}

type Curl struct {
	headers *types.Headers
	flags   *types.Flags
	preset  presets.Preset
	binary  string
	isValid bool
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP Error %d: %s", e.StatusCode, e.Status)
}

func New(options ...Option) *Curl {
	curl := &Curl{
		preset:  presets.Default(),
		headers: types.NewHeaders(),
		flags:   types.NewFlags(types.Flag("include", true)),
	}

	curl.Set(options...)
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

func (curl *Curl) Request(url string) (*http.Response, error) {
	if !curl.isValid {
		if err := curl.Validate(); err != nil {
			return nil, err
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
		return nil, fmt.Errorf("error running curl.\n %s", stderr.String())
	}

	// Get the full output from stdout
	output := stdout.String()

	// Handle potential duplicate HTTP status lines caused by proxy
	lines := strings.Split(output, "\n")
	for i := 1; i < len(lines) && i < 5; i++ { // Limit check to first 5 lines (hardcoded number, potentially can be any humber of extra headers from proxy)
		if strings.HasPrefix(lines[i], "HTTP") {
			// Found a second HTTP status line, remove everything before it
			lines = lines[i:]
			break
		}
	}
	output = strings.Join(lines, "\n")

	// Split the output into headers and body (separated by \r\n\r\n)
	parts := strings.SplitN(output, "\r\n\r\n", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("unable to extract status code from response")
	}

	// The first part contains headers, the second part contains the body
	respHeaders := parts[0]
	respBody := parts[1]

	// Get the status code from the first line of headers
	statusLine := strings.Split(output, "\n")[0]
	statusParts := strings.Split(statusLine, " ")
	if len(statusParts) < 2 {
		return nil, fmt.Errorf("unable to extract status code")
	}
	statusCode, err := getStatusCode(statusParts[1])
	if err != nil {
		return nil, err
	}
	if statusCode >= 400 {
		return nil, &HTTPError{
			StatusCode: statusCode,
			Status:     http.StatusText(statusCode),
		}
	}

	return &http.Response{
		StatusCode: statusCode,
		Header:     parseHeaders(respHeaders),
		Body:       io.NopCloser(bytes.NewReader([]byte(respBody))),
	}, nil
}

func getStatusCode(statusCode string) (int, error) {
	code, err := strconv.Atoi(statusCode)
	if err != nil {
		return 0, fmt.Errorf("invalid status code format: %s", statusCode)
	}
	return code, nil
}

func parseHeaders(headers string) http.Header {
	result := make(http.Header)
	lines := strings.Split(headers, "\r\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result.Add(key, value)
		}
	}
	return result
}

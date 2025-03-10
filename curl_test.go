package curl

import (
	"bytes"
	"errors"
	"testing"
)

func TestCurlLogic(t *testing.T) {
	// Create a new curl instance
	c := New()

	// Test Case 1: Request to an invalid URL (should return a CurlError)
	t.Run("Invalid URL", func(t *testing.T) {
		_, _, _, err := c.Request("http://nonexistent-domain.example")

		// Check if the error is a CurlError
		if !IsCurlError(err) {
			t.Errorf("Expected a CurlError for an invalid URL, got: %v", err)
			return
		}

		// Verify the exit code and error message
		var curlErr *Error
		if errors.As(err, &curlErr) {
			expectedExitCode := 6 // Exit code for "Couldn't resolve host"
			if curlErr.ExitCode != expectedExitCode {
				t.Errorf("Expected exit code %d, got %d", expectedExitCode, curlErr.ExitCode)
			}

			expectedMessage := "Curl Error. Couldn't resolve host. (6)"
			if curlErr.Error() != expectedMessage {
				t.Errorf("Expected error message '%s', got '%s'", expectedMessage, curlErr.Error())
			}
		} else {
			t.Errorf("Expected error to be of type *curl.Error, got: %T", err)
		}
	})

	// Test Case 2: Request to a valid URL but with a timeout (should return a CurlError)
	t.Run("Timeout", func(t *testing.T) {
		// Set a very short timeout to force a timeout error
		c.Set(Flag("max-time", "1")) // Timeout after 1 second

		_, _, _, err := c.Request("http://httpbin.org/delay/5") // This endpoint delays the response by 5 seconds

		// Check if the error is a CurlError
		if !IsCurlError(err) {
			t.Errorf("Expected a CurlError for a timeout, got: %v", err)
			return
		}

		// Verify the exit code and error message
		var curlErr *Error
		if errors.As(err, &curlErr) {
			expectedExitCode := 28 // Exit code for "Operation timeout"
			if curlErr.ExitCode != expectedExitCode {
				t.Errorf("Expected exit code %d, got %d", expectedExitCode, curlErr.ExitCode)
			}

			expectedMessage := "Curl Error. Operation timeout. (28)"
			if curlErr.Error() != expectedMessage {
				t.Errorf("Expected error message '%s', got '%s'", expectedMessage, curlErr.Error())
			}
		} else {
			t.Errorf("Expected error to be of type *curl.Error, got: %T", err)
		}
	})
}

func TestExtractAllResponses(t *testing.T) {
	// Sample output with multiple HTTP responses using \r\n as line endings
	output := []byte("HTTP/1.0 200 Connection established\r\n" +
		"HTTP/2 302 \r\n" +
		"server: nginx\r\n" +
		"date: Sun, 02 Feb 2025 17:27:30 GMT\r\n" +
		"content-type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		"HTTP/1.0 200 Connection established\r\n" +
		"HTTP/2 301 \r\n" +
		"location: https://www.google.com/\r\n" +
		"content-type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		"HTTP/1.0 200 Connection established\r\n" +
		"HTTP/2 200 \r\n" +
		"date: Sun, 02 Feb 2025 17:27:38 GMT\r\n" +
		"expires: -1\r\n" +
		"cache-control: private, max-age=0\r\n" +
		"content-type: text/html; charset=ISO-8859-1\r\n" +
		"\r\n" +
		"<!doctype html><html itemscope=\"\" itemtype=\"http://schema.org/WebPage\" lang=\"en-GB\">")

	// Call extractAllResponses
	headers, lastBody, err := extractAllResponses(output)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Test Case 1: Verify the number of headers
	expectedHeaderCount := 3
	if len(headers) != expectedHeaderCount {
		t.Errorf("Expected %d headers, got %d", expectedHeaderCount, len(headers))
	}

	// Test Case 2: Verify the content of the headers
	expectedHeaders := [][]byte{
		[]byte("HTTP/1.0 200 Connection established\r\n" +
			"HTTP/2 302 \r\n" +
			"server: nginx\r\n" +
			"date: Sun, 02 Feb 2025 17:27:30 GMT\r\n" +
			"content-type: text/html; charset=UTF-8"),
		[]byte("HTTP/1.0 200 Connection established\r\n" +
			"HTTP/2 301 \r\n" +
			"location: https://www.google.com/\r\n" +
			"content-type: text/html; charset=UTF-8"),
		[]byte("HTTP/1.0 200 Connection established\r\n" +
			"HTTP/2 200 \r\n" +
			"date: Sun, 02 Feb 2025 17:27:38 GMT\r\n" +
			"expires: -1\r\n" +
			"cache-control: private, max-age=0\r\n" +
			"content-type: text/html; charset=ISO-8859-1"),
	}
	for i, header := range headers {
		if !bytes.Equal(header, expectedHeaders[i]) {
			t.Errorf("Header mismatch at index %d\nExpected:\n%s\nGot:\n%s", i, expectedHeaders[i], header)
		}
	}

	// Test Case 3: Verify the last body
	expectedLastBody := []byte("<!doctype html><html itemscope=\"\" itemtype=\"http://schema.org/WebPage\" lang=\"en-GB\">")
	if !bytes.Equal(lastBody, expectedLastBody) {
		t.Errorf("Expected last body '%s', got '%s'", expectedLastBody, lastBody)
	}
}

func TestExtractAllResponses_EmptyInput(t *testing.T) {
	// Test Case: Empty input
	output := []byte("")
	_, _, err := extractAllResponses(output)

	// Check if an error is returned
	if err == nil {
		t.Errorf("Expected an error for empty input, got nil")
	}

	// Verify the error message
	expectedError := "unable to extract valid HTTP responses"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestExtractAllResponses_NoValidHeaders(t *testing.T) {
	// Test Case: Input without valid HTTP headers
	output := []byte("This is not a valid HTTP response")
	_, _, err := extractAllResponses(output)

	// Check if an error is returned
	if err == nil {
		t.Errorf("Expected an error for invalid input, got nil")
	}

	// Verify the error message
	expectedError := "unable to extract valid HTTP responses"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

package curl

import (
	"errors"
	"fmt"
)

type HTTPError struct {
	StatusCode int
	Status     string
}

type Error struct {
	ExitCode int
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP Error. %s (%d)", e.Status, e.StatusCode)
}

func (e *Error) Error() string {
	if desc, ok := curlExitCodes[e.ExitCode]; ok {
		return fmt.Sprintf("Curl Error. %s (%d)", desc, e.ExitCode)
	}
	return fmt.Sprintf("Curl Error. Unknown error (%d)", e.ExitCode)
}

func IsHttpError(err error) bool {
	var e *HTTPError
	return errors.As(err, &e)
}

func IsCurlError(err error) bool {
	var e *Error
	return errors.As(err, &e)
}

var curlExitCodes = map[int]string{
	0:  "Success.",
	1:  "Unsupported protocol.",
	2:  "Failed to initialize.",
	3:  "URL malformed.",
	4:  "A feature or option passed to curl is not supported.",
	5:  "Couldn't resolve proxy.",
	6:  "Couldn't resolve host.",
	7:  "Failed to connect to host.",
	8:  "Weird server reply.",
	22: "HTTP page not retrieved.",
	23: "Write error.",
	26: "Read error.",
	27: "Out of memory.",
	28: "Operation timeout.",
	33: "HTTP range error.",
	34: "HTTP post error.",
	35: "SSL connect error.",
	36: "Bad download resume.",
	37: "File couldn't read file.",
	38: "LDAP cannot bind.",
	39: "LDAP search failed.",
	41: "Function not found.",
	42: "Aborted by callback.",
	43: "Internal error.",
	45: "Interface error.",
	47: "Too many redirects.",
	48: "Unknown option specified.",
	49: "Malformed telnet option.",
	51: "The peer's SSL certificate or SSH MD5 fingerprint was not OK.",
	52: "The server didn't reply anything.",
	53: "SSL crypto engine not found.",
	54: "Cannot set SSL crypto engine as default.",
	55: "Failed sending network data.",
	56: "Failure in receiving network data.",
	58: "Problem with the local certificate.",
	59: "Couldn't use specified SSL cipher.",
	60: "Peer certificate cannot be authenticated.",
	61: "Unrecognized transfer encoding.",
	62: "Invalid LDAP URL.",
	63: "Maximum file size exceeded.",
	77: "Problem reading the SSL CA cert.",
	78: "The resource referenced in the URL does not exist.",
	79: "An unspecified error occurred during the SSH session.",
	80: "Failed to shut down the SSL connection.",
	82: "Could not load CRL file.",
	83: "Issuer check failed.",
	85: "Mismatch of RTSP CSeq numbers.",
	86: "Mismatch of RTSP Session Identifiers.",
	88: "Chunk callback reported error.",
	89: "Stream error in the chunked-encoding.",
	90: "Not present in this build.",
	91: "Invalid TLS options were provided.",
	92: "A requested feature, protocol, or option was not found built-in.",
	93: "Out of critical resources.",
	94: "Interface failed to initialize properly.",
	95: "An authentication function returned an error.",
	96: "A problem was detected in the HTTP/3 layer.",
}

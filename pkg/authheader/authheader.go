package authheader

import (
	"errors"
	"strings"
)

const (
	basic  = "basic"
	bearer = "bearer"
)

var (
	ErrInvalidAuthHeader = errors.New("invalid authorization header")
)

func valid(header string, ofType string) (string, error) {
	header = strings.TrimSpace(header)
	h, t := len(header), len(ofType)
	// Must have at least a space and a character after, e.g. `Bearer x`.
	if h < t+2 || !strings.EqualFold(header[0:t], ofType) {
		return "", ErrInvalidAuthHeader
	}
	return header[t+1:], nil
}

// Basic returns the basic authorization token from the header.
func Basic(header string) (string, error) {
	return valid(header, basic)
}

// Bearer returns the bearer authorization token from the header.
func Bearer(header string) (string, error) {
	return valid(header, bearer)
}

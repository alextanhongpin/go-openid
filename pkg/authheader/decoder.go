package authheader

import (
	"encoding/base64"
	"errors"
	"strings"
)

// DecodeBase64 decodes a base64 authorization header into separate fields.
func DecodeBase64(data string) (string, string, error) {
	dec, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return "", "", err
	}
	cred := string(dec)
	if idx := strings.Index(cred, sep); idx > 0 {
		return cred[0:idx], cred[idx+1:], nil
	}
	return "", "", errors.New("index out of range")
}

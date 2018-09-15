package authheader

import (
	"encoding/base64"
	"strings"
)

// DecodeBase64 decodes a base64 authorization header into separate fields.
func DecodeBase64(data string) (string, string) {
	dec, _ := base64.URLEncoding.DecodeString(data)
	cred := string(dec)
	if idx := strings.Index(cred, sep); idx > 0 {
		return cred[0:idx], cred[idx+1:]
	}
	return "", ""
}

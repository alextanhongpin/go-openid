package oidc

import (
	"encoding/base64"
	"strings"
)

// EncodeBasicAuth encodes the username and password into url-safe base64 string.
func EncodeBasicAuth(username, password string) string {
	data := username + ":" + password
	return base64.URLEncoding.EncodeToString([]byte(data))
}

// DecodeBasicAuth decodes a base64 string into the corresponding username and password.
func DecodeBasicAuth(data string) (username, password string) {
	dec, _ := base64.URLEncoding.DecodeString(data)
	cred := string(dec)
	if idx := strings.Index(cred, ":"); idx > 0 {
		return cred[0:idx], cred[idx+1:]
	}
	return "", ""
}

// DecodeClientAuth decodes an authorization header into the corresponding client id and secret.
func DecodeClientAuth(data string) (clientID, clientSecret string) {
	dec, _ := base64.URLEncoding.DecodeString(data)
	cred := string(dec)
	if idx := strings.Index(cred, ":"); idx > 0 {
		return cred[0:idx], cred[idx+1:]
	}
	return "", ""
}

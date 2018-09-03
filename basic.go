package oidc

import (
	"encoding/base64"
	"strings"
)

func EncodeBasicAuth(username, password string) string {
	data := username + ":" + password
	return base64.URLEncoding.EncodeToString([]byte(data))
}

func DecodeBasicAuth(data string) (username, password string) {
	dec, _ := base64.URLEncoding.DecodeString(data)
	cred := string(dec)
	if idx := strings.Index(cred, ":"); idx > 0 {
		return cred[0:idx], cred[idx+1:]
	}
	return "", ""
}

func DecodeClientAuth(data string) (clientID, clientSecret string) {
	dec, _ := base64.URLEncoding.DecodeString(data)
	cred := string(dec)
	if idx := strings.Index(cred, ":"); idx > 0 {
		return cred[0:idx], cred[idx+1:]
	}
	return "", ""
}

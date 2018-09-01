package oidc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeBasicAuth(t *testing.T) {
	assert := assert.New(t)

	var (
		username = "john"
		password = "password"
	)
	enc := EncodeBasicAuth(username, password)
	u, p := DecodeBasicAuth(enc)
	assert.Equal(username, u, "should match the given username")
	assert.Equal(password, p, "should match the given password")
}

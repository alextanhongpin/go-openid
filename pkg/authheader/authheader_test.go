package authheader_test

import (
	"testing"
	"testing/quick"

	"github.com/alextanhongpin/go-openid/pkg/authheader"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		header string
		err    error
		token  string
	}{
		{"Basic abc", nil, "abc"},
		{"basic ab c", nil, "ab c"},
		{"basic abc", nil, "abc"},
		{"basic ", authheader.ErrInvalidAuthHeader, ""},
		{"basic", authheader.ErrInvalidAuthHeader, ""},
		{"", authheader.ErrInvalidAuthHeader, ""},
		{"b", authheader.ErrInvalidAuthHeader, ""},
	}

	for _, tt := range tests {
		token, err := authheader.Basic(tt.header)
		assert.Equal(tt.err, err, "should match the error")
		assert.Equal(tt.token, token, "should match the token")
	}
}

func TestBearer(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		header string
		err    error
		token  string
	}{
		{"Bearer abc", nil, "abc"},
		{"Bearer ab c", nil, "ab c"},
		{"bearer abc", nil, "abc"},
		{"Bearer ", authheader.ErrInvalidAuthHeader, ""},
		{"Bearerr", authheader.ErrInvalidAuthHeader, ""},
		{"", authheader.ErrInvalidAuthHeader, ""},
		{"b", authheader.ErrInvalidAuthHeader, ""},
	}

	for _, tt := range tests {
		token, err := authheader.Bearer(tt.header)
		assert.Equal(tt.err, err, "should match the error")
		assert.Equal(tt.token, token, "should match the token")
	}
}

func TestPanic(t *testing.T) {
	f := func(s string) bool {
		token, err := authheader.Bearer(s)
		return (len(token) > 0 && err == nil) || (len(token) == 0 && err != nil)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Fatal(err)
	}
}

func TestEncodeDecodeBasicAuth(t *testing.T) {
	assert := assert.New(t)

	var (
		username = "john"
		password = "password"
	)

	enc := authheader.EncodeBase64(username, password)
	u, p, err := authheader.DecodeBase64(enc)
	assert.Nil(err)
	assert.Equal(username, u, "should match the given username")
	assert.Equal(password, p, "should match the given password")
}

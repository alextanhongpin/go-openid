package session_test

import (
	"testing"

	"github.com/alextanhongpin/go-openid/pkg/session"
	"github.com/stretchr/testify/assert"
)

func TestCookie(t *testing.T) {
	assert := assert.New(t)
	cookie := session.NewCookie("xyz")

	var (
		name     = "id"
		value    = "xyz"
		maxAge   = 1200
		httpOnly = true
	)

	assert.Equal(name, cookie.Name, "should have a generic name for the cookie")
	assert.Equal(value, cookie.Value, "should have the set value")
	assert.Equal(maxAge, cookie.MaxAge, "should have default max-age of 20 minutes")
	assert.Equal(httpOnly, cookie.HttpOnly, "should default to true for httpOnly")
}

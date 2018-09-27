package session_test

import (
	"testing"
	"time"

	"github.com/alextanhongpin/go-openid/pkg/session"
	"github.com/stretchr/testify/assert"
)

func TestCookie(t *testing.T) {
	assert := assert.New(t)
	cookie := session.NewCookie("xyz", time.Unix(0, 10))

	var (
		name     = "id"
		value    = "xyz"
		maxAge   = 1200
		httpOnly = true
		path     = "/"
		expires  = time.Unix(0, 10).Add(20 * time.Minute)
	)

	assert.Equal(name, cookie.Name, "should have a generic name for the cookie")
	assert.Equal(value, cookie.Value, "should have the set value")
	assert.Equal(maxAge, cookie.MaxAge, "should have default max-age of 20 minutes")
	assert.Equal(expires, cookie.Expires, "should have expire time set to 20 minutes from now")
	assert.Equal(path, cookie.Path, "should have default path set to /")
	assert.Equal(httpOnly, cookie.HttpOnly, "should default to true for httpOnly")
}

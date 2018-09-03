package crypto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	assert := assert.New(t)

	c := New("signature")

	var (
		aud = "john"
		sub = "doe"
		iss = "go-openid"
		dur = 10 * time.Minute
	)

	tok, err := c.NewJWT(aud, sub, iss, dur)
	assert.Nil(err)

	claims, err := c.ParseJWT(tok)
	assert.Nil(err)

	assert.Equal(aud, claims.Audience, "should match the audience")
	assert.Equal(sub, claims.Subject, "should match the subject")
	assert.Equal(iss, claims.Issuer, "should match the issuer")
}

package oidc_test

import (
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"

	"github.com/alextanhongpin/go-openid"
)

func TestIDToken(t *testing.T) {
	assert := assert.New(t)

	var (
		str string

		iss   = "test"
		key   = []byte("secret")
		name  = "john"
		nonce = "xyz"
	)

	claims := &oidc.IDToken{
		StandardClaims: jwt.StandardClaims{
			Issuer: iss,
		},
		Profile: &oidc.Profile{
			Name: name,
		},
		Nonce: nonce,
	}
	t.Run("signing key method 1", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		var err error
		str, err = token.SignedString(key)
		assert.Nil(err)
	})

	t.Run("signing key method", func(t *testing.T) {
		ss, err := claims.SignHS256(key)
		assert.Nil(err)
		assert.Equal(ss, str)
	})

	t.Run("validating key", func(t *testing.T) {
		token, err := jwt.ParseWithClaims(str, &oidc.IDToken{}, func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})
		assert.Nil(err)

		claims, ok := token.Claims.(*oidc.IDToken)

		assert.True(ok && token.Valid)
		assert.Equal(iss, claims.StandardClaims.Issuer, "issuer should be equal")
		assert.Equal(name, claims.Profile.Name, "name should be equal")
		assert.Equal(nonce, claims.Nonce, "nonce should be equal")
	})
	t.Run("validating key method", func(t *testing.T) {
		o := oidc.NewIDToken()
		err := o.ParseHS256(str, key)
		assert.Nil(err)

		assert.Equal(iss, o.StandardClaims.Issuer, "issuer should be equal")
		assert.Equal(name, o.Profile.Name, "name should be equal")
		assert.Equal(nonce, o.Nonce, "nonce should be equal")
	})
	t.Run("sign verify key method", func(t *testing.T) {
		o := oidc.NewIDToken()
		o.StandardClaims.Audience = "1"
		o.StandardClaims.Subject = "100"
		ss, err := o.SignHS256(key)
		assert.Nil(err)

		oo := oidc.NewIDToken()
		err = oo.ParseHS256(ss, key)
		assert.Nil(err)

		assert.Equal(o.StandardClaims.Audience, oo.StandardClaims.Audience, "audience should be equal")
		assert.Equal(o.StandardClaims.Subject, oo.StandardClaims.Subject, "subject should be equal")
		assert.Equal(&o, &oo, "should have different address")
	})
}

package token

import jwt "github.com/dgrijalva/jwt-go"

// Option represents the fields to override.
type Option func(*Claims)

// Claims is an alias to the jwt standard claims.
type Claims = jwt.StandardClaims

// NewClaims return the new claims with the options to override the default
// claims.
func NewClaims(claims *Claims, opts ...Option) *Claims {
	for _, o := range opts {
		o(claims)
	}
	return claims
}

// Audience represents the claims audience.
func Audience(aud string) Option {
	return func(c *Claims) {
		c.Audience = aud
	}
}

// ExpiresAt represents the claims expire at time in unix.
func ExpiresAt(exp int64) Option {
	return func(c *Claims) {
		c.ExpiresAt = exp
	}
}

// ID represents the claims identifier.
func ID(id string) Option {
	return func(c *Claims) {
		c.Id = id
	}
}

// IssuedAt represents the claims issued at time.
func IssuedAt(iat int64) Option {
	return func(c *Claims) {
		c.IssuedAt = iat
	}
}

// Issuer represents the claims issuer.
func Issuer(iss string) Option {
	return func(c *Claims) {
		c.Issuer = iss
	}
}

// NotBefore indicates the duration before the token is valid.
func NotBefore(nbf int64) Option {
	return func(c *Claims) {
		c.NotBefore = nbf
	}
}

// Subject represents the claims subject.
func Subject(sub string) Option {
	return func(c *Claims) {
		c.Subject = sub
	}
}

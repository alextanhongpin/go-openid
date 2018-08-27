package oidc

import (
	"time"
)

const CodeTTL = 10 * time.Minute

// Code represents the authorization code.
type Code struct {
	Code      string
	CreatedAt time.Time
	TTL       time.Duration
}

// NewCode returns a new code with the default TTL.
func NewCode(code string) *Code {
	return &Code{
		Code:      code,
		CreatedAt: time.Now(),
		TTL:       CodeTTL,
	}
}

// Expired returns if the code has reached pass the expiration limit.
func (c *Code) Expired() bool {
	return time.Since(c.CreatedAt) > c.TTL
}

package openid

import (
	"time"
)

// Code represents the authorization code
type Code struct {
	Code      string
	CreatedAt time.Time
	TTL       time.Duration
}

// Expired returns if the code has reached pass the expiration limit
func (c *Code) Expired() bool {
	return time.Since(c.CreatedAt) > c.TTL
}

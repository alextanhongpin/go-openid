



package model

import (
	"context"
	"time"

	"github.com/alextanhongpin/go-openid"
)

// Core represents the core model for the OpenID Connect Specification.
type Core interface {
	// NewCode generates a one-time authorization code that will expire
	// after the set TTL.
	NewCode() string

	// ValidateAuthnRequest validates the required fields for the
	// authentication request.
	ValidateAuthnRequest(req *openid.AuthenticationRequest) error

	// ValidateAuthnClient validates the request payload with the client
	// info from the database.
	ValidateAuthnClient(req *openid.AuthenticationRequest) error

	// ValidateAuthnUser validates the request payload with the user info
	// from the database.
	ValidateAuthnUser(ctx context.Context, req *openid.AuthenticationRequest) error

	ValidateClientAuthHeader(authorization string) (*openid.Client, error)

	// TODO: rename to NewToken/NewIDToken.
	ProvideToken(userID string, duration time.Duration) (string, error)
	ProvideIDToken(userID string) (string, error)
}

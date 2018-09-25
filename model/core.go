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
	ValidateAuthnRequest(req *oidc.AuthenticationRequest) error

	// ValidateAuthnClient validates the request payload with the client
	// info from the database.
	ValidateAuthnClient(req *oidc.AuthenticationRequest) error

	// ValidateAuthnUser validates the request payload with the user info
	// from the database.
	ValidateAuthnUser(ctx context.Context, req *oidc.AuthenticationRequest) error

	ValidateClientAuthHeader(authorization string) (*oidc.Client, error)

	ProvideToken(userID string, duration time.Duration) (string, error)
	ProvideIDToken(userID string) (string, error)
}

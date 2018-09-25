package service

import (
	"context"

	"github.com/alextanhongpin/go-openid"
)

// Core represents the core service for the OpenID Connect Specification.
type Core interface {
	Authenticate(context.Context, *oidc.AuthenticationRequest) (*oidc.AuthenticationResponse, error)
}

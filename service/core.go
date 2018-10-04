




package service

import (
	"context"

	"github.com/alextanhongpin/go-openid"
)

// Core represents the core service for the OpenID Connect Specification.
type Core interface {
	PreAuthenticate(*openid.AuthenticationRequest) error
	Authenticate(context.Context, *openid.AuthenticationRequest) (*openid.AuthenticationResponse, error)
	Token(context.Context, *openid.AccessTokenRequest) (*openid.AccessTokenResponse, error)
}

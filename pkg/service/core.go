package service

import "github.com/alextanhongpin/go-openid"

// Core represents the core service for the OpenID Connect Specification.
type Core interface {
	Authenticate(*oidc.AuthenticationRequest) (*oidc.AuthenticationResponse, error)
}

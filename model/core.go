package model

import "github.com/alextanhongpin/go-openid"

// Core represents the core model for the OpenID Connect Specification.
type Core interface {
	GetClient(id string) (*oidc.Client, error)
	GetUser(id string) (*oidc.User, error)
	NewCode() string
}

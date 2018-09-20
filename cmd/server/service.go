package main

import (
	"github.com/alextanhongpin/go-openid/internal/user"
	"github.com/alextanhongpin/go-openid/pkg/service"
)

// Service represents the interface for the services available for OpenID Connect Protocol.
type Service interface {
	service.User
	// Authorize(context.Context, *oidc.AuthorizationRequest) (*oidc.AuthorizationResponse, error)
	// Token(context.Context, *oidc.AccessTokenRequest) (*oidc.AccessTokenResponse, error)
	// RegisterClient(context.Context, *oidc.Client) (*oidc.Client, error)
	// ReadClient(context.Context, string) (*oidc.Client, error)
	// UserInfo(context.Context, string) (*oidc.IDToken, error)
}

type serviceImpl struct {
	service.User
}

// NewService returns a new service.
func NewService() *serviceImpl {
	u := user.NewService()
	return &serviceImpl{User: u}
}

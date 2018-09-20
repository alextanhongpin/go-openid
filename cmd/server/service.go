package main

import (
	"context"

	oidc "github.com/alextanhongpin/go-openid"
)

// Service represents the interface for the services available for OpenID Connect Protocol.
type Service interface {
	Authorize(context.Context, *oidc.AuthorizationRequest) (*oidc.AuthorizationResponse, error)
	Token(context.Context, *oidc.AccessTokenRequest) (*oidc.AccessTokenResponse, error)
	RegisterClient(context.Context, *oidc.Client) (*oidc.Client, error)
	ReadClient(context.Context, string) (*oidc.Client, error)
	UserInfo(context.Context, string) (*oidc.IDToken, error)
}

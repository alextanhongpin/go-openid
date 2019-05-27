package usecase

import (
	"context"

	openid "github.com/alextanhongpin/go-openid"
)

type Client interface {
	// Register a new client and return a client with client id and client
	// secret.
	Register(ctx context.Context, req RegisterClientRequest) (*RegisterClientResponse, error)
}

type (
	RegisterClientRequest struct {
		Data openid.Client
	}

	RegisterClientResponse struct {
		ClientID     string
		ClientSecret string
	}
)

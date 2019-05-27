package client

import (
	"context"

	"github.com/alextanhongpin/go-openid/domain/client"
	"github.com/alextanhongpin/go-openid/usecase"
)

type UseCase struct {
	clients client.Repository
}

func (u *UseCase) Register(ctx context.Context, req usecase.RegisterClientRequest) (usecase.RegisterClientResponse, error) {
	// Get the current time from the context to make testing easier.
	// ts, ok := contextKeyTimestamp
	// Override.
	success, err := u.clients.Create(req.Data)
	return nil, err
}

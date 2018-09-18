// +build wireinject

package client

import (
	"github.com/google/go-cloud/wire"

	"github.com/alextanhongpin/go-openid/pkg/schema"
)

// NewService returns a new client service.
func NewService(map[string]schema.Validator) *clientServiceImpl {
	panic(wire.Build(ClientServiceSet))
}

// NewModel returns a new client model.
func NewModel(map[string]schema.Validator) *clientModelImpl {
	panic(wire.Build(ClientModelSet))
}

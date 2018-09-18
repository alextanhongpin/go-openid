// +build wireinject

package client

import (
	"github.com/google/go-cloud/wire"

	"github.com/alextanhongpin/go-openid/pkg/schema"
)

// NewService returns a new client service.
func NewService(schema.Validators) *clientServiceImpl {
	panic(wire.Build(ClientServiceSet))
}

// NewModel returns a new client model.
func NewModel(schema.Validators) *clientModelImpl {
	panic(wire.Build(ClientModelSet))
}

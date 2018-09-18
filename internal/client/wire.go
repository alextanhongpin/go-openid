// +build wireinject

package client

import (
	"github.com/google/go-cloud/wire"

	"github.com/alextanhongpin/go-openid/pkg/schema"
)

func NewService(map[string]schema.Validator) *clientServiceImpl {
	panic(wire.Build(ClientSet))
}

func NewSimplerService(map[string]schema.Validator) *clientServiceImpl {
	panic(wire.Build(ClientMegaSet))
}

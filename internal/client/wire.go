// +build wireinject

package client

import (
	"github.com/google/go-cloud/wire"
)

// NewService returns a new client service.
func NewService() (*clientServiceImpl, error) {
	panic(wire.Build(ClientServiceSet))
}

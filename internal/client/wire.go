// +build wireinject

package client

import (
	database "github.com/alextanhongpin/go-openid/internal/database"
	"github.com/alextanhongpin/go-openid/model"
	"github.com/alextanhongpin/go-openid/repository"
	"github.com/google/go-cloud/wire"
)

// $ wire
var clientServiceSet = wire.NewSet(
	provideRepository,
	provideModel,
	NewService,
)

// New returns a new client service.
func New() *Service {
	panic(wire.Build(clientServiceSet))
}

func provideRepository() repository.Client {
	return database.NewClientKV()
}

func provideModel() model.Client {
	return NewModel()
}

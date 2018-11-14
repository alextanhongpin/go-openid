// +build wireinject

package user

import (
	"github.com/google/go-cloud/wire"

	"github.com/alextanhongpin/go-openid/internal/database"
	"github.com/alextanhongpin/go-openid/repository"
)

var serviceSet = wire.NewSet(
	provideRepository,
	wire.Bind(new(repository.User), new(database.UserKV)),
	// provideModel,
	// wire.Bind(new(model.User), new(Model)),
	provideService,
)

// New returns a new  service.
func New() *Service {
	panic(wire.Build(serviceSet))
}

func provideRepository() *database.UserKV {
	return database.NewUserKV()
}

func provideService(repository repository.User) *Service {
	return &Service{
		repository: repository,
	}
}

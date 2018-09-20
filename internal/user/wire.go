// +build wireinject

package user

import (
	"github.com/google/go-cloud/wire"

	"github.com/alextanhongpin/go-openid/internal/database"
	"github.com/alextanhongpin/go-openid/pkg/model"
	"github.com/alextanhongpin/go-openid/pkg/repository"
)

var userServiceSet = wire.NewSet(
	provideRepository,
	wire.Bind(new(repository.User), new(database.UserKV)),
	provideModel,
	wire.Bind(new(model.User), new(userModelImpl)),
	provideService,
)

// NewService returns a new user service.
func NewService() *userServiceImpl {
	panic(wire.Build(userServiceSet))
}

func provideRepository() *database.UserKV {
	return database.NewUserKV()
}

func provideModel(repo repository.User) *userModelImpl {
	return &userModelImpl{repository: repo}
}

func provideService(model model.User) *userServiceImpl {
	decorateValidator := &userValidatorImpl{model: model}
	return &userServiceImpl{model: decorateValidator}
}

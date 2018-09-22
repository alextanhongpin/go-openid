// +build wireinject

package user

import (
	"github.com/google/go-cloud/wire"

	"github.com/alextanhongpin/go-openid/internal/database"
	"github.com/alextanhongpin/go-openid/pkg/model"
	"github.com/alextanhongpin/go-openid/pkg/repository"
)

var ServiceSet = wire.NewSet(
	provideRepository,
	wire.Bind(new(repository.User), new(database.UserKV)),
	provideModel,
	wire.Bind(new(model.User), new(modelImpl)),
	provideService,
)

// NewService returns a new  service.
func NewService() *serviceImpl {
	panic(wire.Build(ServiceSet))
}

func provideRepository() *database.UserKV {
	return database.NewUserKV()
}

func provideModel(repo repository.User) *modelImpl {
	return &modelImpl{repository: repo}
}

func provideService(model model.User) *serviceImpl {
	decorateValidator := &validatorImpl{model: model}
	return &serviceImpl{model: decorateValidator}
}

// +build wireinject

package core

import (
	"github.com/alextanhongpin/go-openid/internal/database"
	"github.com/alextanhongpin/go-openid/model"
	"github.com/alextanhongpin/go-openid/repository"
	"github.com/google/go-cloud/wire"
)

var serviceSet = wire.NewSet(
	provideClientRepository,
	wire.Bind(new(repository.Client), new(database.ClientKV)),
	provideCodeRepository,
	wire.Bind(new(repository.Code), new(database.CodeKV)),
	provideUserRepository,
	wire.Bind(new(repository.User), new(database.UserKV)),
	provideModel,
	wire.Bind(new(model.Core), new(modelImpl)),
	provideService,
)

func NewService() *serviceImpl {
	panic(wire.Build(serviceSet))
}

func provideClientRepository() *database.ClientKV {
	return database.NewClientKV()
}

func provideCodeRepository() *database.CodeKV {
	return database.NewCodeKV()
}

func provideUserRepository() *database.UserKV {
	return database.NewUserKV()
}

func provideModel(code repository.Code, client repository.Client, user repository.User) *modelImpl {
	return &modelImpl{code, client, user}
}

func provideService(model model.Core) *serviceImpl {
	return &serviceImpl{model}
}

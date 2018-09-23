// +build wireinject

package core

import (
	"github.com/alextanhongpin/go-openid/internal/database"
	"github.com/alextanhongpin/go-openid/pkg/model"
	"github.com/alextanhongpin/go-openid/pkg/repository"
	"github.com/google/go-cloud/wire"
)

var serviceSet = wire.NewSet(
	provideClientRepository,
	wire.Bind(new(repository.Client), new(database.ClientKV)),
	provideCodeRepository,
	wire.Bind(new(repository.Code), new(database.CodeKV)),
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

func provideModel(code repository.Code, client repository.Client) *modelImpl {
	return &modelImpl{code, client}
}

func provideService(model model.Core) *serviceImpl {
	return &serviceImpl{model}
}

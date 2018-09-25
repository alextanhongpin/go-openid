// +build wireinject

package client

import (
	database "github.com/alextanhongpin/go-openid/internal/database"
	schema "github.com/alextanhongpin/go-openid/pkg/schema"
	"github.com/alextanhongpin/go-openid/repository"
	"github.com/google/go-cloud/wire"
)

var clientServiceSet = wire.NewSet(
	provideRepository,
	provideModel,
	provideClientValidator,
	provideClientResponseValidator,
	provideValidator,
	provideService,
)

// NewService returns a new client service.
func NewService() (*serviceImpl, error) {
	panic(wire.Build(clientServiceSet))
}

func provideRepository() repository.Client {
	return database.NewClientKV()
}

func provideModel(repo repository.Client) *modelImpl {
	return NewModel(repo)
}

func provideClientValidator() (*schema.Client, error) {
	return schema.NewClientValidator()
}

func provideClientResponseValidator() (*schema.ClientResponse, error) {
	return schema.NewClientResponseValidator()
}

func provideValidator(model *modelImpl, client *schema.Client, clientResponse *schema.ClientResponse) *validatorImpl {
	return &validatorImpl{
		model:          model,
		client:         client,
		clientResponse: clientResponse,
	}
}

func provideService(model *validatorImpl) *serviceImpl {
	return &serviceImpl{model}
}

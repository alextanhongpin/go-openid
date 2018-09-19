// +build wireinject

package client

import (
	database "github.com/alextanhongpin/go-openid/internal/database"
	"github.com/alextanhongpin/go-openid/pkg/repository"
	schema "github.com/alextanhongpin/go-openid/pkg/schema"
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
func NewService() (*clientServiceImpl, error) {
	panic(wire.Build(clientServiceSet))
}

func provideRepository() repository.Client {
	return database.NewClientKV()
}

func provideModel(repo repository.Client) *clientModelImpl {
	return NewClientModelImpl(repo)
}

func provideClientValidator() (*schema.Client, error) {
	return schema.NewClientValidator()
}

func provideClientResponseValidator() (*schema.ClientResponse, error) {
	return schema.NewClientResponseValidator()
}

func provideValidator(model *clientModelImpl, client *schema.Client, clientResponse *schema.ClientResponse) *clientValidatorImpl {
	return &clientValidatorImpl{
		model:          model,
		client:         client,
		clientResponse: clientResponse,
	}
}

func provideService(model *clientValidatorImpl) *clientServiceImpl {
	return NewClientServiceImpl(model)
}

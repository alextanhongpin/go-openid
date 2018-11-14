// +build wireinject

package client

import (
	"github.com/google/go-cloud/wire"

	database "github.com/alextanhongpin/go-openid/internal/database"
	"github.com/alextanhongpin/go-openid/model"
	"github.com/alextanhongpin/go-openid/repository"
)

// $ wire
var clientServiceSet = wire.NewSet(
	provideRepository,
	provideValidator,
	provideModel,
	NewService,
)

// New returns a new client service.
func New() (*Service, error) {
	panic(wire.Build(clientServiceSet))
}

func provideRepository() repository.Client {
	return database.NewClientKV()
}

func provideModel(validator *Validator) model.Client {
	return NewModel(validator)
}

func provideValidator() (*Validator, error) {
	return NewValidator()
}

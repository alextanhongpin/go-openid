package testdata

import (
	"github.com/stretchr/testify/mock"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/database"
)

type clientValidator struct {
	mock.Mock
}

func NewClientValidator() *clientValidator {
	return &clientValidator{}
}

func (c *clientValidator) Validate(client *oidc.Client) error {
	args := c.Called(client)
	return args.Error(0)
}

type clientRepository struct {
	*database.ClientKV
	mock.Mock
}

func NewClientRepository() *clientRepository {
	return &clientRepository{
		ClientKV: database.NewClientKV(),
	}
}

func (c *clientRepository) GenerateClientCredentials() (string, string) {
	args := c.Called()
	return args.String(0), args.String(1)
}

package testdata

import (
	"github.com/stretchr/testify/mock"

	"github.com/alextanhongpin/go-openid/internal/database"
)

type clientRepository struct {
	*database.ClientKV
	mock.Mock
}

func NewClientRepository() *clientRepository {
	return &clientRepository{
		ClientKV: database.NewClientKV(),
	}
}

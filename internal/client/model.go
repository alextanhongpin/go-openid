package client

import (
	"errors"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/repository"
)

type modelImpl struct {
	repository repository.Client
}

// NewModel returns a new client model implementation.
func NewModel(r repository.Client) *modelImpl {
	return &modelImpl{r}
}

// New returns a new client with client id and client secret.
func (c *modelImpl) New(client *openid.Client) (*openid.Client, error) {
	return NewClient(client)
}

// Save stores the new, non-existing client into the database.
func (c *modelImpl) Save(client *openid.Client) error {
	if exist := c.repository.Has(client.ClientID); exist {
		return errors.New("client already exist")
	}
	return c.repository.Put(client.ClientID, client)
}

func (c *modelImpl) CheckExist(exist bool) error {
	if exist {
		return errors.New("client already exist")
	}
	return nil
}

// Read returns a client by client_id from the repository.
func (c *modelImpl) Read(clientID string) (*openid.Client, error) {
	return c.repository.Get(clientID)
}

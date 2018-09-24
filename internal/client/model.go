package client

import (
	"errors"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/repository"
)

type modelImpl struct {
	repository repository.Client
}

// NewModel returns a new client model implementation.
func NewModel(r repository.Client) *modelImpl {
	return &modelImpl{r}
}

// New returns a new client with client id and client secret.
func (c *modelImpl) New(client *oidc.Client) (*oidc.Client, error) {
	return NewClient(client)
}

// Save stores the new, non-existing client into the database.
func (c *modelImpl) Save(client *oidc.Client) error {
	if exist := c.repository.Has(client.ClientID); exist {
		return errors.New("client already exist")
	}
	return c.repository.Put(client.ClientID, client)
}

// Read returns a client by client_id from the repository.
func (c *modelImpl) Read(clientID string) (*oidc.Client, error) {
	return c.repository.Get(clientID)
}

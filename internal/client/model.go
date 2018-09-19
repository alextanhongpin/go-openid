package client

import (
	"errors"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/repository"
)

type clientModelImpl struct {
	repository repository.Client
}

// NewClientModelImpl returns a new client model implementation.
func NewClientModelImpl(r repository.Client) *clientModelImpl {
	return &clientModelImpl{r}
}

// New returns a new client with client id and client secret.
func (c *clientModelImpl) New(client *oidc.Client) (*oidc.Client, error) {
	return NewClient(client)
}

// Save stores the new, non-existing client into the database.
func (c *clientModelImpl) Save(client *oidc.Client) error {
	if exist := c.repository.Has(client.ClientID); exist {
		return errors.New("client already exist")
	}
	return c.repository.Put(client.ClientID, client)
}

// Read returns a client by client_id from the repository.
func (c *clientModelImpl) Read(clientID string) (*oidc.Client, error) {
	client, exist := c.repository.Get(clientID)
	if !exist {
		return nil, errors.New("client does not exist")
	}
	return client, nil
}

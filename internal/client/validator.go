





package client

import (
	"errors"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/model"
	schema "github.com/alextanhongpin/go-openid/pkg/schema"
)

type validatorImpl struct {
	client         *schema.Client
	clientResponse *schema.ClientResponse
	model          model.Client
}

// New will validate the new request.
func (c *validatorImpl) New(client *openid.Client) (*openid.Client, error) {
	if client == nil {
		return nil, errors.New("arguments cannot be nil")
	}
	_, err := c.client.Validate(client)
	if err != nil {
		return nil, err
	}
	return c.model.New(client)
}

// Save will validate the save request.
func (c *validatorImpl) Save(client *openid.Client) error {
	if client == nil {
		return errors.New("arguments cannot be nil")
	}
	_, err := c.clientResponse.Validate(client)
	if err != nil {
		return err
	}
	return c.model.Save(client)
}

// Read will validate the read request.
func (c *validatorImpl) Read(clientID string) (*openid.Client, error) {

	if clientID == "" {
		return nil, errors.New("client_id cannot be empty")
	}
	return c.model.Read(clientID)
}

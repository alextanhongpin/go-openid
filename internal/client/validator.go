package client

import (
	"errors"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/model"
	"github.com/alextanhongpin/go-openid/pkg/schema"
)

type clientValidatorImpl struct {
	client         *schema.Client
	clientResponse *schema.ClientResponse
	model          model.Client
}

// New will validate the new request.
func (c *clientValidatorImpl) New(client *oidc.Client) (*oidc.Client, error) {
	_, err := c.client.Validate(client)
	if err != nil {
		return nil, err
	}
	return c.model.New(client)
}

// Save will validate the save request.
func (c *clientValidatorImpl) Save(client *oidc.Client) error {
	_, err := c.clientResponse.Validate(client)
	if err != nil {
		return err
	}
	return c.model.Save(client)
}

// Read will validate the read request.
func (c *clientValidatorImpl) Read(clientID string) (*oidc.Client, error) {
	if clientID == "" {
		return nil, errors.New("client_id cannot be empty")
	}
	return c.model.Read(clientID)
}

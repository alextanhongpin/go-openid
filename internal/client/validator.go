package client

import (
	"errors"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/schema"
)

// TODO: Example validation using decorator pattern.
type clientValidatorImpl struct {
	client         *schema.Client
	clientResponse *schema.ClientResponse
	model          *clientModelImpl
}

// func NewClientValidatorImpl(
//         model *clientModelImpl,
//         client *schema.Client,
//         clientResponse *schema.ClientResponse,
// ) *clientValidatorImpl {
//         return &clientValidatorImpl{
//                 model:          model,
//                 client:         client,
//                 clientResponse: clientResponse,
//         }
// }

func (c *clientValidatorImpl) New(client *oidc.Client) (*oidc.Client, error) {
	_, err := c.client.Validate(client)
	if err != nil {
		return nil, err
	}
	return c.model.New(client)
}

func (c *clientValidatorImpl) Save(client *oidc.Client) error {
	_, err := c.clientResponse.Validate(client)
	if err != nil {
		return err
	}
	return c.model.Save(client)
}

func (c *clientValidatorImpl) Read(clientID string) (*oidc.Client, error) {
	if clientID == "" {
		return nil, errors.New("client_id cannot be empty")
	}
	return c.model.Read(clientID)
}

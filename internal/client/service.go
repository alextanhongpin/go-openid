package client

import (
	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/model"
)

type serviceImpl struct {
	// We need to give this a name in case of collision between the same
	// name, and we want to avoid using the CamelCase client to represent
	// model naming.
	model model.Client
}

// Register performs client registration which will return a new client with
// client id and client secret.
func (c *serviceImpl) Register(client *openid.Client) (*openid.Client, error) {
	newClient, err := c.model.New(client)
	if err != nil {
		return nil, err
	}

	return newClient, c.model.Save(newClient)
}

// Read returns a client by client id or error if the client is not found.
func (c *serviceImpl) Read(clientID string) (*openid.Client, error) {
	return c.model.Read(clientID)
}

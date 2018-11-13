package client

import (
	"errors"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/model"
	"github.com/alextanhongpin/go-openid/repository"
)

type Service struct {
	// We need to give this a name in case of collision between the same
	// name, and we want to avoid using the CamelCase client to represent
	// model naming.
	model model.Client

	repo repository.Client
}

// NewService returns a new Client service.
func NewService(model model.Client, repo repository.Client) *Service {
	return &Service{model, repo}
}

// Register performs client registration which will return a new client with
// client id and client secret.
func (s *Service) Register(client *openid.Client) (*openid.Client, error) {
	if client == nil {
		return nil, errors.New("arguments cannot be nil")
	}
	m := s.model
	modifiers := []ClientModifier{
		ClientValidation(m.ValidateClient),
		ClientID(m.GenerateClientID),
		ClientSecret(m.GenerateClientSecret),
		RegistrationAccessToken(m.GenerateRegistrationAccessToken),
		ClientIDIssuedAt(m.GenerateClientIDIssuedAt),
		ClientSecretExpiresAt(m.GenerateClientSecretExpiresAt),
		RegistrationClientURI(m.GenerateRegistrationClientURI),
		ClientResponseValidation(m.ValidateClientResponse),
	}
	err := apply(client, modifiers...)
	if err != nil {
		return nil, err
	}

	// Repository does not modify the model. It modifies the repository only.
	exist := s.repo.Has(client.ClientID)
	if exist {
		return nil, errors.New("client already exist")
	}
	err = s.repo.Put(client.ClientID, client)
	return client, err
}

// Read returns a client by client id or error if the client is not found.
func (c *Service) Read(clientID string) (*openid.Client, error) {
	if clientID == "" {
		return nil, errors.New("client_id cannot be empty")
	}
	return c.repo.Get(clientID)
}

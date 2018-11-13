package client

import openid "github.com/alextanhongpin/go-openid"

// NopService represents a no-op service for testing edge cases or mocking the service.
type NopService struct{}

func (n *NopService) Register(c *openid.Client) (*openid.Client, error) {
	return c, nil
}

func (n *NopService) Read(clientID string) (*openid.Client, error) {
	return nil, nil
}

type NopModel struct{}

func (n *NopModel) GenerateClientID() string {
	return ""
}

func (n *NopModel) GenerateClientSecret() (string, error) {
	return "", nil
}

func (n *NopModel) GenerateRegistrationAccessToken(string) (string, error) {
	return "", nil
}

func (n *NopModel) GenerateClientIDIssuedAt() int64 {
	return 0
}
func (n *NopModel) GenerateClientSecretExpiresAt() int64 {
	return 0
}

func (n *NopModel) GenerateRegistrationClientURI() string {
	return ""
}

func (n *NopModel) ValidateClient(o *openid.Client) error {
	return nil
}

func (n *NopModel) ValidateClientResponse(o *openid.Client) error {
	return nil
}

type NopRepository struct{}

func (n *NopRepository) Get(id string) (*openid.Client, error) {
	return nil, nil
}

func (n *NopRepository) Put(id string, client *openid.Client) error {
	return nil
}

func (n *NopRepository) Has(id string) bool {
	return false
}

func (n *NopRepository) List(limit int) []*openid.Client {
	return []*openid.Client{}
}

func (n *NopRepository) GetByCredentials(clientID, clientSecret string) (*openid.Client, error) {
	return nil, nil
}

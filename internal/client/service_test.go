package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	openid "github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/client"
)

func TestServiceWithNop(t *testing.T) {
	assert := assert.New(t)

	model := new(client.NopModel)
	repo := new(client.NopRepository)
	service := client.NewService(model, repo)

	client := new(openid.Client)
	client, err := service.Register(client)
	assert.Nil(err)
	// TODO: Test if it populates with default values.
	assert.NotNil(client)
	assert.Equal("", client.ClientID)
	assert.Equal("", client.ClientSecret)
	assert.Equal("", client.RegistrationAccessToken)
	assert.Equal(int64(0), client.ClientIDIssuedAt)
	assert.Equal(int64(0), client.ClientSecretExpiresAt)
	assert.Equal("", client.RegistrationClientURI)

	client, err = service.Read("")
	assert.Equal("client_id cannot be empty", err.Error())
	assert.Nil(client)

}

type repoMock struct {
	client.NopRepository
}

func (r *repoMock) Has(clientID string) bool {
	return true
}

func TestService_ClientExist(t *testing.T) {
	assert := assert.New(t)

	model := new(client.NopModel)
	repo := new(repoMock)
	service := client.NewService(model, repo)

	client := new(openid.Client)
	client, err := service.Register(client)
	assert.NotNil(err)
	assert.Equal("client already exist", err.Error())
}

type modelMock struct {
	client.NopModel
	clientID     string
	clientSecret string
}

func (m *modelMock) GenerateClientID() string {
	return m.clientID
}

func (m *modelMock) GenerateClientSecret() (string, error) {
	return m.clientSecret, nil
}

func TestService_ClientIDSet(t *testing.T) {
	assert := assert.New(t)

	var (
		clientID     = "client_id"
		clientSecret = "client_secret"
	)
	model := &modelMock{
		clientID:     clientID,
		clientSecret: clientSecret,
	}
	repo := new(client.NopRepository)
	service := client.NewService(model, repo)

	client := new(openid.Client)
	client, err := service.Register(client)
	assert.Nil(err)
	assert.Equal(clientID, client.ClientID, "should set the correct client id")
	assert.Equal(clientSecret, client.ClientSecret, "should set the correct client secret")
}

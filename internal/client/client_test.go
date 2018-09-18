package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/client"
	"github.com/alextanhongpin/go-openid/pkg/schema"
)

func TestNewClientModel(t *testing.T) {
	assert := assert.New(t)

	model := client.NewModel(newValidator(t))

	t.Run("register with empty request", func(t *testing.T) {
		client := new(oidc.Client)
		_, err := model.New(client)
		assert.Equal("redirect_uris is required", err.Error(), "should validate the only required field")
	})

	t.Run("register with additional field client_id", func(t *testing.T) {
		client := new(oidc.Client)

		// Minimum required field is redirect_uris, which is fulfilled.
		client.RedirectURIs = oidc.RedirectURIs([]string{"https://server.example.com/cb"})

		// Attempt to inject client_id to override system.
		client.ClientID = "xyz"

		_, err := model.New(client)
		assert.Equal("Additional property client_id is not allowed", err.Error(), "should return error indicating redirect_uris is required")
	})

	t.Run("register and save with only redirect_uris", func(t *testing.T) {
		client := new(oidc.Client)
		client.RedirectURIs = oidc.RedirectURIs([]string{"https://server.example.com/cb"})
		newClient, err := model.New(client)
		assert.Nil(err)

		assert.Equal("", newClient.ApplicationType, "should have default application type of web")
		assert.Equal(0, len(newClient.GrantTypes), "should have default grant_type of authorization_code")
		assert.Equal(1, len(newClient.RedirectURIs), "should match the redirect_uris provided")
		assert.Equal("https://server.example.com/cb", newClient.RedirectURIs[0], "should match the redirect_uris provided")
		assert.Equal("", newClient.RequestObjectEncryptionEnc, "should have default value")
		assert.Equal("", newClient.UserinfoEncryptedResponseEnc, "should have defaut value")

		assert.True(newClient.ClientID != "", "should return a generated client_id")
		assert.True(newClient.ClientSecret != "", "should return a generated client_secret")
		assert.True(newClient.ClientSecretExpiresAt == 0, "should set client_secret_expires_at to infinity")
		assert.True(newClient.RegistrationAccessToken != "", "should return a default registration token")
		assert.True(newClient.RegistrationClientURI == "https://server.example.com/c2id/clients", "should have the correct registration_client_uri")

		err = model.Save(newClient)
		assert.Nil(err)
	})

	t.Run("register and save with default values", func(t *testing.T) {
		client := oidc.NewClient()
		client.RedirectURIs = oidc.RedirectURIs([]string{"https://server.example.com/cb"})

		newClient, err := model.New(client)
		assert.Nil(err)

		var (
			applicationType              = "web"
			grantType                    = "authorization_code"
			redirectURI                  = "https://server.example.com/cb"
			requestObjectEncryptionEnc   = "A128CBC-HS256"
			userinfoEncryptedResponseEnc = "A128CBC-HS256"
		)

		assert.Equal(applicationType, newClient.ApplicationType, "should have default application type of web")
		assert.Equal(grantType, newClient.GrantTypes[0], "should have default grant_type of authorization_code")
		assert.Equal(redirectURI, newClient.RedirectURIs[0], "should match the redirect_uris provided")
		assert.Equal(requestObjectEncryptionEnc, newClient.RequestObjectEncryptionEnc, "should have default value")
		assert.Equal(userinfoEncryptedResponseEnc, newClient.UserinfoEncryptedResponseEnc, "should have defaut value")

		err = model.Save(newClient)
		assert.Nil(err)

		storedClient, err := model.Read(newClient.ClientID)
		assert.Nil(err)
		assert.NotNil(storedClient)
	})

	t.Run("attempt to save an empty client", func(t *testing.T) {
		client := oidc.NewClient()
		err := model.Save(client)
		assert.NotNil(err)

		assert.Equal("client_id is required", err.Error(), "should return the first validation error")
	})

	t.Run("attempt to read a non-existing client", func(t *testing.T) {
		_, err := model.Read("")
		assert.Equal("client does not exist", err.Error())
	})
}

// -- helpers

func die(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func newValidator(t *testing.T) map[string]schema.Validator {
	// Setup validation.
	validateNewRequest, err := schema.NewClientValidator()
	die(t, err)

	// Setup validation.
	validateSaveRequest, err := schema.NewClientResponseValidator()
	die(t, err)

	return map[string]schema.Validator{
		"New":  validateNewRequest,
		"Save": validateSaveRequest,
	}
}

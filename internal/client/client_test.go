package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/client"
)

func TestNewClientService(t *testing.T) {
	assert := assert.New(t)

	service, err := client.New()
	assert.Nil(err)

	t.Run("register with nil request", func(t *testing.T) {
		_, err := service.Register(nil)
		assert.Equal("arguments cannot be nil", err.Error(), "should handle nil arguments")
	})

	t.Run("register with empty request", func(t *testing.T) {
		client := new(openid.Client)
		_, err := service.Register(client)
		assert.Equal("redirect_uris is required", err.Error(), "should validate the only required field")
	})

	t.Run("register with additional field client_id", func(t *testing.T) {
		client := new(openid.Client)

		// Minimum required field is redirect_uris, which is fulfilled.
		client.RedirectURIs = openid.RedirectURIs([]string{"https://server.example.com/cb"})

		// Attempt to inject client_id to override system.
		client.ClientID = "xyz"

		_, err := service.Register(client)
		assert.Equal("Additional property client_id is not allowed", err.Error(), "should return error indicating additional property is not allowed")
	})

	t.Run("register and save with only redirect_uris", func(t *testing.T) {
		client := new(openid.Client)
		client.RedirectURIs = openid.RedirectURIs([]string{"https://server.example.com/cb"})
		newClient, err := service.Register(client)
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

	})

	t.Run("register and save with default values", func(t *testing.T) {
		client := openid.NewClient()
		client.RedirectURIs = openid.RedirectURIs([]string{"https://server.example.com/cb"})

		newClient, err := service.Register(client)
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

		storedClient, err := service.Read(newClient.ClientID)
		assert.Nil(err)
		assert.NotNil(storedClient)
	})

	t.Run("attempt to read a non-existing client", func(t *testing.T) {
		_, err := service.Read("")
		assert.Equal("client_id cannot be empty", err.Error())
	})
}

func TestClientRegistration(t *testing.T) {
	assert := assert.New(t)

	body := []byte(`{
		"application_type": "web",
		"redirect_uris": [
			"https://client.example.org/callback",
			"https://client.example.org/callback2"
		],
		"client_name": "My Example",
		"client_name#ja-Jpan-JP": "クライアント名",
		"logo_uri": "https://client.example.org/logo.png",
		"subject_type": "pairwise",
		"sector_identifier_uri": "https://other.example.net/file_of_redirect_uris.json",
		"token_endpoint_auth_method": "client_secret_basic",
		"jwks_uri": "https://client.example.org/my_public_keys.jwks",
		"userinfo_encrypted_response_alg": "RSA1_5",
		"userinfo_encrypted_response_enc": "A128CBC-HS256",
		"contacts": [
			"ve7jtb@example.org",
			"mary@example.org"
		],
		"request_uri": [
			"https://client.example.org/rf.txt#qpXaRLh_n93TTR9F252ValdatUQvQiJi5BDub2BeznA"
		]
	}`)

	c := openid.NewClient()
	err := c.UnmarshalJSON(body)
	assert.Nil(err)

	service, err := client.New()
	assert.Nil(err)

	newClient, err := service.Register(c)
	assert.Nil(err)
	assert.NotNil(newClient)

	assert.Equal(c.ClientName, newClient.ClientName, "should set the client_name")
	assert.Equal(c.ApplicationType, newClient.ApplicationType, "should set the application_type")
}

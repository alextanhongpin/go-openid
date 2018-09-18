package api_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/api"
	"github.com/alextanhongpin/go-openid/pkg/api/testdata"
	"github.com/alextanhongpin/go-openid/pkg/schema"
)

func TestClientError(t *testing.T) {
	assert := assert.New(t)

	var (
		clientID     = "abc"
		clientSecret = "xyz"
		clientName   = "app"
	)

	req := &oidc.Client{
		ApplicationType:              "",
		ClientName:                   clientName,
		ClientURI:                    "",
		Contacts:                     []string(nil),
		DefaultAcrValues:             "",
		DefaultMaxAge:                0,
		GrantTypes:                   []string(nil),
		IDTokenEncryptedResponseAlg:  "",
		IDTokenEncryptedResponseEnc:  "",
		IDTokenSignedResponseAlg:     "",
		InitiateLoginURI:             "",
		Jwks:                         "",
		JwksURI:                      "",
		LogoURI:                      "",
		PolicyURI:                    "",
		RedirectURIs:                 oidc.RedirectURIs([]string{"http://client.example.com"}),
		RequestObjectEncryptionAlg:   "",
		RequestObjectEncryptionEnc:   "",
		RequestObjectSigningAlg:      "",
		RequestURIs:                  []string(nil),
		RequireAuthTime:              0,
		ResponseTypes:                []string(nil),
		SectorIdentifierURI:          "",
		SubjectType:                  "",
		TokenEndpointAuthMethod:      "",
		TokenEndpointAuthSigningAlg:  "",
		TosURI:                       "",
		UserinfoEncryptedResponseAlg: "",
		UserinfoEncryptedResponseEnc: "",
		UserinfoSignedResponseAlg:    "",
		RegistrationClientURI:        "",
	}

	res := req.Copy()
	res.ClientID = clientID
	res.ClientSecret = clientSecret
	res.ClientSecretExpiresAt = 0
	res.RegistrationAccessToken = ""
	res.RegistrationClientURI = ""

	registerClient := func(client *oidc.Client) (*oidc.Client, error) {
		factory := api.NewClient
		repository := testdata.NewClientRepository()
		requestValidator, err := schema.NewClientValidator()
		if err != nil {
			return nil, err
		}
		responseValidator, err := schema.NewClientRegistrationResponseValidator()
		if err != nil {
			return nil, err
		}
		return api.RegisterClient(client, factory, repository, requestValidator, responseValidator)
	}

	t.Run("test register client", func(t *testing.T) {
		res, err := registerClient(req)
		assert.Nil(err)
		log.Println(res, err)
		assert.True(res.ClientID != "", "should return client id")
		assert.True(res.ClientSecret != "", "should return client secret")
		// assert.Equal(clientSecret, res.ClientSecret)
		// assert.Equal(clientName, res.ClientName)
		// assert.Equal(iat, res.ClientIDIssuedAt)
		// assert.Equal(exp, res.ClientSecretExpiresAt)
	})

	// t.Run("test double-registration", func(t *testing.T) {
	//         _, err := api.RegisterClient(req, clock, repository, validator)
	//         assert.NotNil(err, "should not allow double registration")
	//
	//         assert.Equal("client already exist", err.Error(), "should return the correct error message")
	// })
}

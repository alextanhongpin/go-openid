package api_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/api"
	"github.com/alextanhongpin/go-openid/pkg/api/testdata"
	"github.com/alextanhongpin/go-openid/schema"
)

func TestClientError(t *testing.T) {
	assert := assert.New(t)

	var (
		clientID     = "abc"
		clientSecret = "xyz"
		clientName   = "app"
		iat          = int64(0)
		exp          = int64(1000)
	)

	repository := testdata.NewClientRepository()
	repository.On("GenerateClientCredentials").Return(clientID, clientSecret)

	validator, err := schema.NewClientValidator()
	assert.Nil(err)
	assert.NotNil(validator)

	clock := func() (int64, int64) {
		return iat, exp
	}

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
		ClientID:                     clientID,
		ClientIDIssuedAt:             iat,
		ClientSecret:                 clientSecret,
		ClientSecretExpiresAt:        exp,
		RegistrationAccessToken:      "",
		RegistrationClientURI:        "",
	}

	t.Run("test register client", func(t *testing.T) {
		res, err := api.RegisterClient(req, clock, repository, validator)
		assert.Nil(err)
		log.Println("register client", res, err)

		assert.Equal(clientID, res.ClientID)
		assert.Equal(clientSecret, res.ClientSecret)
		assert.Equal(clientName, res.ClientName)
		assert.Equal(iat, res.ClientIDIssuedAt)
		assert.Equal(exp, res.ClientSecretExpiresAt)

		spew.Dump(res)
		resValidator, err := schema.NewClientRegistrationResponseValidator()
		assert.Nil(err)

		result, err := resValidator.Validate(res)
		assert.Nil(err)
		for _, err := range result.Errors() {
			fmt.Printf("- %s\n", err)
		}

	})
	// t.Run("test double-registration", func(t *testing.T) {
	//         _, err := api.RegisterClient(req, clock, repository, validator)
	//         assert.NotNil(err, "should not allow double registration")
	//
	//         assert.Equal("client already exist", err.Error(), "should return the correct error message")
	// })
}

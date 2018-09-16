package api_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/database"
	"github.com/alextanhongpin/go-openid/pkg/api"
	"github.com/alextanhongpin/go-openid/pkg/api/testdata"
)

func TestClientError(t *testing.T) {
	assert := assert.New(t)

	now, _ := time.Parse(time.RFC3339, "2018-09-16 06:58:11.615034 +0000 UTC")
	var (
		helper     = testdata.NewClientHelper(1*time.Minute, now)
		validator  = testdata.NewClientValidator()
		repository = database.NewClientKV()
	)

	helper.On("NewClientID").Return("abc")
	helper.On("NewClientSecret", "abc", 1*time.Minute).Return("xyz")

	client := &oidc.Client{
		ApplicationType:              "",
		ClientName:                   "",
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
		RedirectURIs:                 oidc.RedirectURIs(nil),
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
		ClientID:                     "abc",
		ClientIDIssuedAt:             0,
		ClientSecret:                 "xyz",
		ClientSecretExpiresAt:        0,
		RegistrationAccessToken:      "",
		RegistrationClientURI:        "",
	}
	validator.On("Validate", client).Return(nil)

	req := new(oidc.Client)
	res, err := api.RegisterClient(req, helper, repository, validator)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal("abc", res.ClientID)
	// assert.Equal("empty client", err.Error())
}

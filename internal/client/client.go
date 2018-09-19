package client

import (
	"time"

	"github.com/google/go-cloud/wire"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/database"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
	"github.com/alextanhongpin/go-openid/pkg/model"
	"github.com/alextanhongpin/go-openid/pkg/repository"
	"github.com/alextanhongpin/go-openid/pkg/schema"
)

// ClientServiceSet represents a new model.
var ClientServiceSet = wire.NewSet(
	// Setup database.
	database.NewClientKV,
	wire.Bind(new(repository.Client), new(database.ClientKV)),
	// Setup validators.
	schema.NewClientValidator,
	schema.NewClientResponseValidator,
	NewClientModelImpl,
	clientValidatorImpl{},
	// Setup service.
	wire.Bind(new(model.Client), new(clientValidatorImpl)),
	NewClientServiceImpl,
)

// NewClient returns a new client with the generated client id and secret.
func NewClient(c *oidc.Client) (*oidc.Client, error) {
	client := c.Clone()

	var (
		clientID = crypto.NewXID()
		iat      = time.Now().UTC()
		day      = time.Hour * 24
		exp      = iat.Add(7 * day)
		aud      = "https://server.example.com/c2id/clients"
		iss      = clientID
		sub      = clientID
		key      = []byte("secret")
	)
	// Generate client secret.
	clientSecret, err := crypto.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}

	// Generate access token.
	claims := crypto.NewStandardClaims(aud, sub, iss, iat.Unix(), exp.Unix())
	accessToken, err := crypto.NewJWT(key, claims)
	if err != nil {
		return nil, err
	}

	client.ClientID = clientID
	client.ClientSecret = clientSecret
	client.ClientIDIssuedAt = iat.Unix()
	client.ClientSecretExpiresAt = 0 // Never expire.
	client.RegistrationAccessToken = accessToken
	client.RegistrationClientURI = aud

	return client, nil
}

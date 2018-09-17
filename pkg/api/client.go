package api

import (
	"errors"
	"time"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/database"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
	"github.com/alextanhongpin/go-openid/pkg/schema"
)

// NewClient returns a new client with the generated client id and secret.
func NewClient(c *oidc.Client) (*oidc.Client, error) {
	client := new(oidc.Client)
	*client = *c

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

// RegisterClient creates a new client and stores it into the repository.
func RegisterClient(
	client *oidc.Client,
	clientFactory func(client *oidc.Client) (*oidc.Client, error),
	repository database.ClientRepository,
	requestValidator schema.Validator,
	responseValidator schema.Validator,
) (*oidc.Client, error) {
	if client == nil {
		return nil, errors.New("empty client")
	}

	if _, err := requestValidator.Validate(client); err != nil {
		return nil, err
	}

	newClient, err := clientFactory(client)
	if err != nil {
		return nil, err
	}

	if _, err := responseValidator.Validate(newClient); err != nil {
		return nil, err
	}

	clientID := newClient.ClientID
	if exist := repository.Has(clientID); exist {
		return nil, errors.New("client already exist")
	}

	// Store everything in the database.
	if err := repository.Put(clientID, newClient); err != nil {
		return nil, err
	}

	return newClient, nil
}

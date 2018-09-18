package api

import (
	"errors"
	"time"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/model"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
	"github.com/alextanhongpin/go-openid/pkg/schema"
)

// -- model

type clientModelImpl struct {
	repository    repository.Client
	newValidator  schema.Validator
	saveValidator schema.Validator
}

// type ClientSaveRequest struct {
//         client *oidc.Client
//         validator schema.Validator
// }

func (c *clientModelImpl) Save(client *oidc.Client) error {
	if _, err := c.saveValidator.Validate(client); err != nil {
		return err
	}
	if exist := c.repository.Has(client.ClientID); exist {
		return errors.New("client already exist")
	}
	return c.repository.Put(client.ClientID, client)
}

func (c *clientModelImpl) New(client *oidc.Client) (*oidc.Client, error) {
	if _, err := c.newValidator.Validate(client); err != nil {
		return nil, err
	}
	return NewClient(client)
}

// -- service

type clientServiceImpl struct {
	model.Client
}

func (c *clientServiceImpl) Register(client *oidc.Client) (*oidc.Client, error) {
	newClient, err := c.New(client)
	if err != nil {
		return nil, err
	}
	return newClient, c.Save(newClient)
}

// -- helper

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

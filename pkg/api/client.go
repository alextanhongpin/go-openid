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
	repository repository.Client
	validators map[string]schema.Validator
}

func NewClientModelImpl(r repository.Client, v map[string]schema.Validator) *clientModelImpl {
	return &clientModelImpl{r, v}
}

func (c *clientModelImpl) Save(client *oidc.Client) error {
	if err := c.validateSave(client); err != nil {
		return err
	}
	if exist := c.repository.Has(client.ClientID); exist {
		return errors.New("client already exist")
	}
	return c.repository.Put(client.ClientID, client)
}

func (c *clientModelImpl) New(client *oidc.Client) (*oidc.Client, error) {
	if err := c.validateNew(client); err != nil {
		return err
	}
	return NewClient(client)
}

// -- model validation

// This means more function, but it is a better way than 1) creating dedicated
// structs with embedded validation, like saveClientRequest.Validate(). Besides
// private functions are not included in the interface, and hence reduce the
// code required for mocking. 2) private functions don't pollute the struct,
// and it's better than dependency injection. This layers do not need mocking
// anyway - we want to skip them in the test, or either test them with actual
// implementations.
func validate(key string, validators map[string]schema.Validator, data interface{}) error {
	// If the validator is not present, skip it.
	if v, ok := validators[key]; !ok {
		return nil
	}
	_, err := v.Validate(data)
	return err
}

func (c *clientModelImpl) validateSave(client *oidc.Client) error {
	return validate("Save", c.validators, client)
}

func (c *clientModelImpl) validateNew(client *oidc.Client) error {
	return validate("New", c.validators, client)
}

// -- service

type clientServiceImpl struct {
	// We need to give this a name in case of collision between the same
	// name, and we want to avoid using the CamelCase client to represent
	// model naming.
	model model.Client
}

func NewClientServiceImpl(m model.Client) *clientServiceImpl {
	return &clientServiceImpl{m}
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

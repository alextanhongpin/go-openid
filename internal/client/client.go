package client

import (
	"errors"
	"time"

	"github.com/google/go-cloud/wire"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/database"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
	"github.com/alextanhongpin/go-openid/pkg/model"
	"github.com/alextanhongpin/go-openid/pkg/repository"
	"github.com/alextanhongpin/go-openid/pkg/schema"
)

// ClientModelSet represents a new model.
var ClientModelSet = wire.NewSet(
	database.NewClientKV,
	wire.Bind(new(repository.Client), new(database.ClientKV)),
	NewClientModelImpl,
)

// ClientServiceSet represents a new service.
var ClientServiceSet = wire.NewSet(
	ClientModelSet,
	wire.Bind(new(model.Client), new(clientModelImpl)),
	NewClientServiceImpl,
)

// -- model

type clientModelImpl struct {
	repository repository.Client
	schema.Validators
}

// NewClientModelImpl returns a new client model implementation.
func NewClientModelImpl(r repository.Client, v schema.Validators) *clientModelImpl {
	return &clientModelImpl{r, v}
}

// New returns a new client with client id and client secret.
func (c *clientModelImpl) New(client *oidc.Client) (*oidc.Client, error) {
	if err := c.validateNew(client); err != nil {
		return nil, err
	}
	return NewClient(client)
}

// Save stores the new, non-existing client into the database.
func (c *clientModelImpl) Save(client *oidc.Client) error {
	if err := c.validateSave(client); err != nil {
		return err
	}
	if exist := c.repository.Has(client.ClientID); exist {
		return errors.New("client already exist")
	}
	return c.repository.Put(client.ClientID, client)
}

// Read returns a client by client_id from the repository.
func (c *clientModelImpl) Read(clientID string) (*oidc.Client, error) {
	client, exist := c.repository.Get(clientID)
	if !exist {
		return nil, errors.New("client does not exist")
	}
	return client, nil
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
	v, ok := validators[key]
	if !ok {
		return nil
	}
	_, err := v.Validate(data)
	return err
}

func (c *clientModelImpl) validateSave(client *oidc.Client) error {
	_, err := c.Validate("Save", client, true)
	return err
}

func (c *clientModelImpl) validateNew(client *oidc.Client) error {
	_, err := c.Validate("New", client, true)
	return err
}

// TODO: Example validation using decorator pattern.
// type clientValidatorImpl struct {
//         save  schema.Validator
//         new   schema.Validator
//         model model.Client
// }
//
//
// func (c *clientValidatorImpl) New(client *oidc.Client) (*oidc.Client, error) {
//         _, err := c.new.Validate(client)
//         if err != nil {
//                 return nil, err
//         }
//         return c.model.New(client)
// }
//
// func (c *clientValidatorImpl) Save(client *oidc.Client) error {
//         _, err := c.save.Validate(client)
//         if err != nil {
//                 return err
//         }
//         return c.model.Save(client)
// }
//
// func (c *clientValidatorImpl) Read(clientID string) (*oidc.Client, error) {
//         if clientID == "" {
//                 return nil, errors.New("client_id cannot be empty")
//         }
//         return c.model.Read(clientID)
// }

// -- service

type clientServiceImpl struct {
	// We need to give this a name in case of collision between the same
	// name, and we want to avoid using the CamelCase client to represent
	// model naming.
	model model.Client
}

// NewClientServiceImpl returns a new client service implementation.
func NewClientServiceImpl(m model.Client) *clientServiceImpl {
	return &clientServiceImpl{m}
}

// Register performs client registration which will return a new client with
// client id and client secret.
func (c *clientServiceImpl) Register(client *oidc.Client) (*oidc.Client, error) {
	newClient, err := c.model.New(client)
	if err != nil {
		return nil, err
	}

	return newClient, c.model.Save(newClient)
}

// Read returns a client by client id or error if the client is not found.
func (c *clientServiceImpl) Read(clientID string) (*oidc.Client, error) {
	return c.model.Read(clientID)
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

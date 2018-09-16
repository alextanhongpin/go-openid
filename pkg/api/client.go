package api

import (
	"errors"
	"log"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/database"
	"github.com/alextanhongpin/go-openid/schema"
)

func newClient(c *oidc.Client, clientID, clientSecret string, iat, exp int64) *oidc.Client {
	client := new(oidc.Client)
	*client = *c

	// Overwrite values set by the client.
	client.ClientID = clientID
	client.ClientSecret = clientSecret
	client.ClientIDIssuedAt = iat
	client.ClientSecretExpiresAt = exp
	client.RegistrationAccessToken = ""
	client.RegistrationClientURI = ""
	return client
}

// helper ClientHelper,
// type ClientHelper interface {
//         NewFrozenTime() (iat, exp int64)
// }

func RegisterClient(
	client *oidc.Client,
	clock func() (iat, exp int64),
	repo database.ClientRepository,
	validator schema.Validator,
) (*oidc.Client, error) {
	if client == nil {
		return nil, errors.New("empty client")
	}

	clientID, clientSecret := repo.GenerateClientCredentials()

	// Generate issued at and expiration time externally; for testability.
	iat, exp := clock()
	newClient := newClient(client, clientID, clientSecret, iat, exp)

	log.Println("ndw client id now", newClient, validator)

	// TODO: additionalProperties in gojsonschema is not working.
	// Validate request with schema before storing.
	if _, err := validator.Validate(newClient); err != nil {
		log.Println("validation error", err)
		return nil, err
	}
	log.Println("ndw client id now", newClient)
	if exist := repo.Has(clientID); exist {
		return nil, errors.New("client already exist")
	}
	// Store everything in the database.
	if err := repo.Put(clientID, newClient); err != nil {
		return nil, err
	}

	return newClient, nil
}

package api

import (
	"time"

	"github.com/alextanhongpin/go-openid"
)

// ClientRepository represents the interface for the client repository.
type ClientRepository interface {
	Has(id string) bool
	GetByID(id string) *oidc.Client
	Put(id string, client *oidc.Client) error
}

type ClientValidator interface {
	Validate(client *oidc.Client) error
}

func RegisterClient(repo ClientRepository, validator ClientValidator, client *oidc.Client) (*oidc.Client, error) {
	// Checking to prevent nil pointer
	client = client.Copy()

	var (
		dur      = 10 * time.Minute
		now      = time.Now().UTC()
		exp      = now.Add(dur)
		clientID = newClientID()
		attempts = 3
	)

	// Attempt to generate a unique client id.
	for i := 0; i < attempts; i++ {
		if ok := repo.Has(clientID); !ok {
			// Break when the id is unique.
			break
		}
		clientID = newClientID()
	}

	// Generate client id, client secret
	clientSecret := newClientSecret(clientID, dur)

	client.ClientID = clientID
	client.ClientSecret = clientSecret
	client.ClientIDIssuedAt = int64(now.Second())
	client.ClientSecretExpiresAt = int64(exp.Second())
	client.RegistrationAccessToken = ""
	client.RegistrationClientURI = ""

	// Validate request with schema before storing.
	if err := validator.Validate(client); err != nil {
		return nil, err
	}

	// Store everything in the database.
	if err := repo.Put(clientID, client); err != nil {
		return nil, err
	}

	return client, nil
}

func newClientID() string {
	return "abc"
}

func newClientSecret(clientID string, duration time.Duration) (clientSecret string) {
	return ""
}

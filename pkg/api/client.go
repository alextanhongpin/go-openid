package api

import (
	"errors"
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

type ClientHelper interface {
	NewDuration() time.Duration
	NewTime() time.Time
	NewClientID() string
	NewClientSecret(clientID string, duration time.Duration) string
}

func RegisterClient(
	client *oidc.Client,
	helper ClientHelper,
	repo ClientRepository,
	validator ClientValidator,
) (*oidc.Client, error) {
	if client == nil {
		return nil, errors.New("empty client")
	}
	// Checking to prevent nil pointer
	client = client.Copy()

	var (
		dur      = helper.NewDuration()
		now      = helper.NewTime()
		exp      = now.Add(dur)
		clientID = helper.NewClientID()
		attempts = 3
	)

	// Attempt to generate a unique client id.
	for i := 0; i < attempts; i++ {
		if ok := repo.Has(clientID); !ok {
			// Break when the id is unique.
			break
		}
		clientID = helper.NewClientID()
	}

	// Generate client id, client secret
	clientSecret := helper.NewClientSecret(clientID, dur)

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

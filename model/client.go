package model

import openid "github.com/alextanhongpin/go-openid"

// Client represents the client model.
type Client interface {
	// NOTE: The following does not fulfil a model - which is performing
	// business logic. Plain CRUD is not business logic.
	// Read(clientID string) (*openid.Client, error)
	// Save(client *openid.Client) error
	// Update
	// Get
	// Delete

	// There are several other idea regarding model - do we need a struct at all? We can probably assume the models are pure functions - they take in an input and returns the one with business logic applied to them. They don't mutate the state though - that is the job of the repository.

	// Our models becomes plainly data access modifiers.
	GenerateClientID() string
	GenerateClientSecret() (string, error)
	GenerateRegistrationAccessToken(string) (string, error)
	GenerateClientIDIssuedAt() int64
	GenerateClientSecretExpiresAt() int64
	GenerateRegistrationClientURI() string
	ValidateClient(o *openid.Client) error
	ValidateClientResponse(o *openid.Client) error
}

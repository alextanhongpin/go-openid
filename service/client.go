package service

import "github.com/alextanhongpin/go-openid"

// Client represents the client service.
type Client interface {
	// Register a new client and return a client with client id and client
	// secret.
	Register(client *openid.Client) (*openid.Client, error)

	// Read returns a client by client id, and returns error if the client
	// is not found.
	Read(clientID string) (*openid.Client, error)
}

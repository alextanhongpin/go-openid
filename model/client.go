package model

import "github.com/alextanhongpin/go-openid"

// Client represents the client model.
type Client interface {
	New(client *oidc.Client) (*oidc.Client, error)
	Read(clientID string) (*oidc.Client, error)
	Save(client *oidc.Client) error
	// Update
	// Get
	// Delete
}

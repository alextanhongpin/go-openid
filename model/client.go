



package model

import "github.com/alextanhongpin/go-openid"

// Client represents the client model.
type Client interface {
	New(client *openid.Client) (*openid.Client, error)
	Read(clientID string) (*openid.Client, error)
	Save(client *openid.Client) error
	// Update
	// Get
	// Delete
}

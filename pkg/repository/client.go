package repository

import "github.com/alextanhongpin/go-openid"

// Client represents the interface for the client repository.
type Client interface {
	Get(id string) (*oidc.Client, error)
	// GetByID(id string) *oidc.Client
	// GetByIDAndSecret(id, secret string) *oidc.Client
	Put(id string, client *oidc.Client) error
	// Delete(name string)
	Has(id string) bool
	List(limit int) []*oidc.Client
}

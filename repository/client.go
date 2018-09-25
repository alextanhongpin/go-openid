package repository

import "github.com/alextanhongpin/go-openid"

// Client represents the interface for the client repository.
type Client interface {
	Get(id string) (*oidc.Client, error)
	Put(id string, client *oidc.Client) error
	Has(id string) bool
	List(limit int) []*oidc.Client
	GetByCredentials(clientID, clientSecret string) (*oidc.Client, error)
}

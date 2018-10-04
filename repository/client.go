




package repository

import "github.com/alextanhongpin/go-openid"

// Client represents the interface for the client repository.
type Client interface {
	Get(id string) (*openid.Client, error)
	Put(id string, client *openid.Client) error
	Has(id string) bool
	List(limit int) []*openid.Client
	GetByCredentials(clientID, clientSecret string) (*openid.Client, error)
}

package repository

import "github.com/alextanhongpin/go-openid"

// User represents the operations for the user repository.
type User interface {
	FindByEmail(email string) (*oidc.User, error)
	Put(id string, user *oidc.User) error
	Get(id string) (*oidc.User, error)
	List(limit int) ([]*oidc.User, error)
}

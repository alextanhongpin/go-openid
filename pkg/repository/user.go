package repository

import "github.com/alextanhongpin/go-openid"

// User represents the operations for the user repository.
type User interface {
	FindByEmail(email string, sanitized bool) (*oidc.User, error)
	Put(id string, user *oidc.User) error
}

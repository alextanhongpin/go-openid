package model

import "github.com/alextanhongpin/go-openid"

// User represents the user model.
type User interface {
	Create(email, password string) (*oidc.User, error)
	FindByEmail(email string) (*oidc.User, error)
}

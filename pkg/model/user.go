package model

import "github.com/alextanhongpin/go-openid"

// User represents the user model.
type User interface {
	Has(email string) error
	Create(email, hash string) error
	FindByEmail(email string) (*oidc.User, error)
}

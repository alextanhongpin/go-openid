package model

import "github.com/alextanhongpin/go-openid"

// User represents the user model.
type User interface {
	Create(email, hash string) error
	FindByEmail(email string, sanitized bool) (*oidc.User, error)
}

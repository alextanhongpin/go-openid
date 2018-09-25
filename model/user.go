package model

import "github.com/alextanhongpin/go-openid"

// User represents the user model.
type User interface {
	// Create returns creates a new user with the given username and
	// password.
	Create(email, password string) (*oidc.User, error)

	// FindByEmail checks if the user with the given email exists.
	FindByEmail(email string) (*oidc.User, error)
}

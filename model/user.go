package model

import "github.com/alextanhongpin/go-openid"

// User represents the user model.
type User interface {
	NewUser(email, password string) (*openid.User, error)
	ValidateEmail(email string) error
	ValidatePassword(password string) error
	ValidateLimit(limit int) error
}

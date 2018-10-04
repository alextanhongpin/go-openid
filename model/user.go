


package model

import "github.com/alextanhongpin/go-openid"

// User represents the user model.
type User interface {
	// Create returns creates a new user with the given username and
	// password.
	Create(email, password string) (*openid.User, error)

	// FindByEmail checks if the user with the given email exists.
	FindByEmail(email string) (*openid.User, error)

	List(limit int) ([]*openid.User, error)

	Get(id string) (*openid.User, error)
}

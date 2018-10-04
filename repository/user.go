package repository

import "github.com/alextanhongpin/go-openid"

// User represents the operations for the user repository.
type User interface {
	FindByEmail(email string) (*openid.User, error)
	Put(id string, user *openid.User) error
	Get(id string) (*openid.User, error)
	List(limit int) ([]*openid.User, error)
}





package service

import "github.com/alextanhongpin/go-openid"

// User represents the user interface.
type User interface {
	Login(email, password string) (*openid.User, error)
	Register(email, password string) (*openid.User, error)
	GetUsers(limit int) ([]*openid.User, error)
	GetUser(id string) (*openid.User, error)
}

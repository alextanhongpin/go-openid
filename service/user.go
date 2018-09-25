package service

import "github.com/alextanhongpin/go-openid"

// User represents the user interface.
type User interface {
	Login(email, password string) (*oidc.User, error)
	Register(email, password string) (*oidc.User, error)
	GetUsers(limit int) ([]*oidc.User, error)
	GetUser(id string) (*oidc.User, error)
}

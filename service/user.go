package service

import "github.com/alextanhongpin/go-openid"

// User represents the user interface.
type User interface {
	Login(email, password string) (*oidc.User, error)
	Register(email, password string) (*oidc.User, error)
}

package usecase

import openid "github.com/alextanhongpin/go-openid"

type UseCase interface {
	Login(email, password string) (*openid.User, error)
	Register(email, password string) (*openid.User, error)
	WithID(id string) (*openid.User, error)
}

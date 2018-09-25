package user

import (
	"errors"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/database"
	model "github.com/alextanhongpin/go-openid/model"
)

type serviceImpl struct {
	model model.User
}

// Register registers a new user with the given email and password.
func (u *serviceImpl) Register(email, password string) (*oidc.User, error) {
	user, err := u.model.FindByEmail(email)
	switch {
	case err == database.ErrEmailDoesNotExist:
		return u.model.Create(email, password)
	case user != nil: // User exist.
		return nil, errors.New("email exist")
	default:
		return nil, err
	}
}

// Login verifies if the given username and password is correct.
func (u *serviceImpl) Login(email, password string) (*oidc.User, error) {
	user, err := u.model.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, user.ComparePassword(password)
}

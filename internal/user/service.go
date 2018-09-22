package user

import (
	"errors"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/database"
	"github.com/alextanhongpin/go-openid/pkg/model"
	"github.com/alextanhongpin/go-openid/pkg/passwd"
)

type serviceImpl struct {
	model model.User
}

// Register registers a new user with the given email and password.
func (u *serviceImpl) Register(email, password string) error {
	user, err := u.model.FindByEmail(email, true)
	switch {
	case err == database.ErrEmailDoesNotExist:
		hash, err := passwd.Hash(password)
		if err != nil {
			return err
		}
		return u.model.Create(email, hash)
	case user != nil: // User exist.
		return errors.New("email exist")
	default:
		return err
	}
}

// Login verifies if the given username and password is correct.
func (u *serviceImpl) Login(email, password string) (*oidc.User, error) {
	user, err := u.model.FindByEmail(email, false)
	if err != nil {
		return nil, err
	}
	return user, passwd.Verify(password, user.HashedPassword)
}

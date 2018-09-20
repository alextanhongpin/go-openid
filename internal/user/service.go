package user

import (
	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/model"
	"github.com/alextanhongpin/go-openid/pkg/passwd"
)

type userServiceImpl struct {
	model model.User
}

// Register registers a new user with the given email and password.
func (u *userServiceImpl) Register(email, password string) error {
	if _, err := u.model.FindByEmail(email); err != nil {
		return err
	}
	hash, err := passwd.Hash(password)
	if err != nil {
		return err
	}
	return u.model.Create(email, hash)
}

// Login verifies if the given username and password is correct.
func (u *userServiceImpl) Login(email, password string) (*oidc.User, error) {
	user, err := u.model.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, passwd.Verify(user.HashedPassword, password)
}

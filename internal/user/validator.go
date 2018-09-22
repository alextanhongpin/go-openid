package user

import (
	"errors"

	"github.com/asaskevich/govalidator"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/model"
)

type validatorImpl struct {
	model model.User
}

func (u *validatorImpl) FindByEmail(email string) (*oidc.User, error) {
	if err := isEmail(email); err != nil {
		return nil, err
	}
	return u.model.FindByEmail(email)
}

func (u *validatorImpl) Create(email, password string) (*oidc.User, error) {
	if len(email) == 0 || len(password) == 0 {
		return nil, errors.New("arguments cannot be empty")
	}
	if err := isEmail(email); err != nil {
		return nil, err
	}
	return u.model.Create(email, password)
}

// -- helpers

func isEmail(email string) error {
	if ok := govalidator.IsEmail(email); !ok {
		return errors.New("invalid email")
	}
	return nil
}

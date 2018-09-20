package user

import (
	"errors"

	"github.com/asaskevich/govalidator"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/model"
)

type userValidatorImpl struct {
	model model.User
}

func (u *userValidatorImpl) FindByEmail(email string) (*oidc.User, error) {
	if err := isEmail(email); err != nil {
		return nil, err
	}
	return u.model.FindByEmail(email)
}

func (u *userValidatorImpl) Create(email, hashedPassword string) error {
	if len(email) == 0 || len(hashedPassword) == 0 {
		return errors.New("arguments cannot be empty")
	}
	if err := isEmail(email); err != nil {
		return err
	}
	return u.model.Create(email, hashedPassword)
}

// -- helpers

func isEmail(email string) error {
	if ok := govalidator.IsEmail(email); !ok {
		return errors.New("invalid email")
	}
	return nil
}

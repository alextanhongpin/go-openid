package user

import (
	"errors"

	"github.com/asaskevich/govalidator"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/model"
)

type validatorImpl struct {
	model model.User
}

func (u *validatorImpl) FindByEmail(email string) (*openid.User, error) {
	if err := isEmail(email); err != nil {
		return nil, err
	}
	return u.model.FindByEmail(email)
}

func (u *validatorImpl) Create(email, password string) (*openid.User, error) {
	if len(email) == 0 || len(password) == 0 {
		return nil, errors.New("arguments cannot be empty")
	}
	if len(password) < 8 {
		return nil, errors.New("password cannot be less than 8 characters")
	}
	if err := isEmail(email); err != nil {
		return nil, err
	}
	return u.model.Create(email, password)
}

func (u *validatorImpl) List(limit int) ([]*openid.User, error) {
	if limit < 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return u.model.List(limit)
}

func (u *validatorImpl) Get(id string) (*openid.User, error) {
	if id == "" {
		return nil, errors.New("user_id cannot be empty")
	}
	return u.model.Get(id)
}

// -- helpers

func isEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if ok := govalidator.IsEmail(email); !ok {
		return errors.New("invalid email")
	}
	return nil
}

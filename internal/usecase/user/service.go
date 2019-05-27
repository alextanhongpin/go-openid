package user

import (
	"errors"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/repository"
	"github.com/alextanhongpin/go-openid/utils"
)

type Service struct {
	repository repository.User
}

func NewService(repository repository.User) *Service {
	return &Service{repository}
}

func validate(email, password string) error {
	var err error
	err = utils.ValidateEmail(email)
	if err != nil {
		return err
	}
	return utils.ValidatePassword(password)
}

// Register registers a new user with the given email and password.
func (s *Service) Register(email, password string) (*openid.User, error) {
	if err := validate(email, password); err != nil {
		return nil, err
	}
	if err := IsRegistered(s.repository, email); err != nil {
		return nil, err
	}
	user := NewUser(email)
	err := user.SetPassword(password)
	if err != nil {
		return nil, err
	}
	err = Register(s.repository, user)
	return user, err
}

// Login verifies if the given username and password is correct.
func (s *Service) Login(email, password string) (*openid.User, error) {
	if err := validate(email, password); err != nil {
		return nil, err
	}
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	err = user.ComparePassword(password)
	return user, err
}

// GetUsers returns a list of paginated users.
func (s *Service) GetUsers(limit int) ([]*openid.User, error) {
	err := utils.ValidateRange(limit, 0, 20)
	if err != nil {
		return nil, err
	}
	return s.repository.List(limit)
}

// GetUser returns a user by id.
func (u *Service) GetUser(id string) (*openid.User, error) {
	if utils.ValidateString(id) {
		return nil, errors.New("id is required")
	}
	return u.repository.Get(id)
}

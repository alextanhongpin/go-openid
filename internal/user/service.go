package user

import (
	"errors"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/database"
	model "github.com/alextanhongpin/go-openid/model"
	"github.com/alextanhongpin/go-openid/repository"
)

type Service struct {
	model      model.User
	repository repository.User
}

func NewService(model model.User, repository repository.User) *Service {
	return &Service{model, repository}
}

// Register registers a new user with the given email and password.
func (s *Service) Register(email, password string) (*openid.User, error) {
	err := s.model.ValidateEmail(email)
	if err != nil {
		return nil, err
	}
	err = s.model.ValidatePassword(password)
	if err != nil {
		return nil, err
	}
	// user, err := s.model.FindByEmail(email)
	user, err := s.repository.FindByEmail(email)
	switch {
	case err == database.ErrEmailDoesNotExist:
		user, err := s.model.NewUser(email, password)
		if err != nil {
			return nil, err
		}
		err = s.repository.Put(user.ID, user)
		return user, err
	case user != nil: // User exist.
		return nil, errors.New("email exist")
	default:
		return nil, err
	}
}

// Login verifies if the given username and password is correct.
func (s *Service) Login(email, password string) (*openid.User, error) {
	err := s.model.ValidateEmail(email)
	if err != nil {
		return nil, err
	}
	err = s.model.ValidatePassword(password)
	if err != nil {
		return nil, err
	}
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, user.ComparePassword(password)
}

// GetUsers returns a list of paginated users.
func (s *Service) GetUsers(limit int) ([]*openid.User, error) {
	err := s.model.ValidateLimit(limit)
	if err != nil {
		return nil, err
	}
	return s.repository.List(limit)
}

// GetUsers returns a list of paginated users.
func (u *Service) GetUser(id string) (*openid.User, error) {
	if IsEmptyString(id) {
		return nil, errors.New("id cannot be empty")
	}
	// return u.model.Get(id)
	return u.repository.Get(id)
}

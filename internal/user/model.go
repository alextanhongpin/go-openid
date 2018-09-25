package user

import (
	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/repository"
)

type modelImpl struct {
	repository repository.User
}

// FindByEmail returns a user by email.
func (m *modelImpl) FindByEmail(email string) (*oidc.User, error) {
	return m.repository.FindByEmail(email)
}

// Create stores the username and hashed password into the storage.
func (m *modelImpl) Create(email, password string) (*oidc.User, error) {
	user := NewUser()
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	user.Email.Email = email

	return user, m.repository.Put(user.ID, user)
}

func (m *modelImpl) List(limit int) ([]*oidc.User, error) {
	return m.repository.List(10)
}

func (m *modelImpl) Get(id string) (*oidc.User, error) {
	return m.repository.Get(id)
}

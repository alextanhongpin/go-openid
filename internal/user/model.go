package user

import (
	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/repository"
)

type modelImpl struct {
	repository repository.User
}

// FindByEmail returns a user by email.
func (m *modelImpl) FindByEmail(email string) (*openid.User, error) {
	return m.repository.FindByEmail(email)
}

// Create stores the username and hashed password into the storage.
func (m *modelImpl) Create(email, password string) (*openid.User, error) {
	user := NewUser()
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	user.Email.Email = email

	return user, m.repository.Put(user.ID, user)
}

// List should return a paginated array of user.
func (m *modelImpl) List(limit int) ([]*openid.User, error) {
	return m.repository.List(10)
}

// Get should return a client by the client id.
func (m *modelImpl) Get(id string) (*openid.User, error) {
	return m.repository.Get(id)
}

package user

import (
	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/repository"
)

type userModelImpl struct {
	repository repository.User
}

// FindByEmail returns a user by email.
func (u *userModelImpl) FindByEmail(email string, sanitized bool) (*oidc.User, error) {
	return u.repository.FindByEmail(email, sanitized)
}

// Create stores the username and hashed password into the storage.
func (u *userModelImpl) Create(email, hashedPassword string) error {
	user := NewUser()

	user.Email.Email = email
	user.HashedPassword = hashedPassword

	return u.repository.Put(user.ID, user)
}

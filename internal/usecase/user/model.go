package user

import (
	"errors"
	"time"

	openid "github.com/alextanhongpin/go-openid"
	database "github.com/alextanhongpin/go-openid/internal/database"
	repository "github.com/alextanhongpin/go-openid/repository"
	"github.com/rs/xid"
)

func Register(repo repository.User, user *openid.User) error {
	return repo.Put(user.ID, user)
}

func IsRegistered(repo repository.User, email string) error {
	user, err := repo.FindByEmail(email)
	if user != nil {
		return errors.New("user already exist")
	}
	if err == database.ErrEmailDoesNotExist {
		return nil
	}
	return err
}

// NewUser creates a new User.
func NewUser(email string) *openid.User {
	return &openid.User{
		ID: xid.New().String(),
		Profile: openid.Profile{
			UpdatedAt: time.Now().UTC().Unix(),
		},
		Email: openid.Email{
			Email: email,
		},
	}
}

package user

import (
	"time"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
)

// NewUser returns a new user with default values.
func NewUser() *openid.User {
	id := crypto.NewXID()

	user := new(openid.User)
	user.ID = id
	user.Profile.UpdatedAt = time.Now().UTC().Unix()

	return user
}

package user

import (
	"time"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
)

// NewUser returns a new user with default values.
func NewUser() *oidc.User {
	id := crypto.NewXID()

	user := oidc.User{}
	user.ID = id
	user.Profile.UpdatedAt = time.Now().UTC().Unix()

	return &user
}

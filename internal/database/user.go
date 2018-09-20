package database

import (
	"errors"
	"sync"

	"github.com/alextanhongpin/go-openid"
)

// Index table.
var emailToIDmap map[string]string
var idToEmailmap map[string]string

// UserKV represents the in-memory store for user.
type UserKV struct {
	sync.RWMutex
	db map[string]*oidc.User
}

// NewUserKV returns a new user key-value store.
func NewUserKV() *UserKV {
	return &UserKV{
		db: make(map[string]*oidc.User),
	}
}

// Get returns a user by id.
// func (u *UserKV) Get(id string) (*oidc.IDToken, bool) {
//         u.RLock()
//         user, exist := u.db[id]
//         u.RUnlock()
//         return user, exist
// }

// Put stores the user in the db by the given id.
func (u *UserKV) Put(id string, user *oidc.User) {
	u.Lock()
	u.db[id] = user
	u.Unlock()
}

// Delete removes the user with the given id.
// func (u *UserKV) Delete(id string) {
//         u.Lock()
//         delete(u.db, id)
//         u.Unlock()
// }

func (u *UserKV) FindByEmail(email string) (*oidc.User, error) {
	u.RLock()
	user, exist := u.db[email]
	u.RUnlock()

	if !exist {
		return nil, errors.New("email does not exist")
	}
	return user, nil
}

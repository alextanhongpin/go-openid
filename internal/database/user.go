package database

import (
	"sync"

	"github.com/alextanhongpin/go-openid"
)

// UserKV represents the in-memory store for user.
type UserKV struct {
	sync.RWMutex
	db map[string]*oidc.StandardClaims
}

// NewUserKV returns a new user key-value store.
func NewUserKV() *UserKV {
	return &UserKV{
		db: make(map[string]*oidc.StandardClaims),
	}
}

// Get returns a user by id.
func (u *UserKV) Get(id string) (*oidc.StandardClaims, bool) {
	u.RLock()
	user, exist := u.db[id]
	u.RUnlock()
	return user, exist
}

// Put stores the user in the db by the given id.
func (u *UserKV) Put(id string, user *oidc.StandardClaims) {
	u.Lock()
	u.db[id] = user
	u.Unlock()
}

// Delete removes the user with the given id.
func (u *UserKV) Delete(id string) {
	u.Lock()
	delete(u.db, id)
	u.Unlock()
}

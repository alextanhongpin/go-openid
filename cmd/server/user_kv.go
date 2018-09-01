package main

import "sync"
import "github.com/alextanhongpin/go-openid"

type User struct {
	ID string `json:"id"`
	*oidc.StandardClaims
}

func NewUser(id string, sc *oidc.StandardClaims) *User {
	return &User{
		ID:             id,
		StandardClaims: sc,
	}
}

// UserKV represents the in-memory store for user.
type UserKV struct {
	sync.RWMutex
	db map[string]*User
}

// NewUserKV returns a new user key-value store.
func NewUserKV() *UserKV {
	return &UserKV{
		db: make(map[string]*User),
	}
}

// Get returns a user by id.
func (u *UserKV) Get(id string) (*User, bool) {
	u.RLock()
	user, exist := u.db[id]
	u.RUnlock()
	return user, exist
}

// Put stores the user in the db by the given id.
func (u *UserKV) Put(id string, user *User) {
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

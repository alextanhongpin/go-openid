package database

import (
	"errors"
	"sync"

	"github.com/alextanhongpin/go-openid"
)

var (
	ErrEmailDoesNotExist = errors.New("email does not exist")
)

// UserKV represents the in-memory store for user.
type UserKV struct {
	sync.RWMutex
	db map[string]*openid.User
	// Maps email to id, and vice-versa.
	idx map[string]string
}

// NewUserKV returns a new user key-value store.
func NewUserKV() *UserKV {
	return &UserKV{
		db:  make(map[string]*openid.User),
		idx: make(map[string]string),
	}
}

// Get returns a user by id.
func (u *UserKV) Get(id string) (*openid.User, error) {
	u.RLock()
	user, exist := u.db[id]
	u.RUnlock()
	if !exist {
		return nil, errors.New("user does not exist")
	}
	return user, nil
}

func (u *UserKV) List(limit int) (users []*openid.User, err error) {
	u.RLock()
	defer u.RUnlock()
	i := 0
	for _, v := range u.db {
		i++
		users = append(users, v)
		if i >= limit {
			break
		}
	}
	return
}

// Put stores the user in the db by the given id.
func (u *UserKV) Put(id string, user *openid.User) error {
	email := user.Email.Email
	u.Lock()
	u.db[id] = user
	u.Unlock()

	// Set indices.
	u.Lock()
	u.idx[email] = id
	u.Unlock()
	return nil
}

// Delete removes the user with the given id.
// func (u *UserKV) Delete(id string) {
//         u.Lock()
//         delete(u.db, id)
//         u.Unlock()
// }

// FindByEmail returns a user by the given email, or error if the email does
// not exist.
func (u *UserKV) FindByEmail(email string) (*openid.User, error) {
	u.RLock()
	id, exist := u.idx[email]
	u.RUnlock()
	if !exist {
		return nil, ErrEmailDoesNotExist
	}

	u.RLock()
	user, exist := u.db[id]
	u.RUnlock()
	if !exist {
		return nil, ErrEmailDoesNotExist
	}
	return user, nil
}

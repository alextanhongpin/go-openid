package main

import "sync"

type User struct{}

// In-memory store for user
type UserKV struct {
	sync.RWMutex
	db map[string]*User
}

func NewUserKV() *UserKV {
	return &UserKV{
		db: make(map[string]*User),
	}
}

func (u *UserKV) Get(id string) (*User, bool) {
	u.RLock()
	user, exist := u.db[id]
	u.RUnlock()
	return user, exist
}

func (u *UserKV) Put(id string, user *User) {
	u.Lock()
	u.db[id] = user
	u.Unlock()
}

func (u *UserKV) Delete(id string) {
	u.Lock()
	delete(u.db, id)
	u.Unlock()
}

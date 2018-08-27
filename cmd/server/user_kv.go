package main

import "sync"

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

func (u *UserKV) Get(id string) (user *User, exist bool) {
	u.RLock()
	user, exist = u.db[id]
	u.RUnlock()
}

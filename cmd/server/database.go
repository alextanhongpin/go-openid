package main

type Database struct {
	Client *ClientKV
	Code   *CodeKV
	User   *UserKV
}

// NewDatabase returns an in-memory database.
func NewDatabase() *Database {
	return &Database{
		Client: NewClientKV(),
		Code:   NewCodeKV(),
		User:   NewUserKV(),
	}
}

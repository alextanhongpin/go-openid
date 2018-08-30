package main

// Database represents the data storage access layer.
type Database struct {
	Client ClientRepository
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

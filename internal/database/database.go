package database

import "github.com/alextanhongpin/go-openid/pkg/repository"

// Database represents the data storage access layer.
type Database struct {
	Client repository.Client
	Code   repository.Code
	User   repository.User
}

// NewInMem returns an in-memory database.
func NewInMem() *Database {
	return &Database{
		Client: NewClientKV(),
		Code:   NewCodeKV(),
		User:   NewUserKV(),
	}
}

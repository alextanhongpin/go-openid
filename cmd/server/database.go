package main

type Database struct {
	Client *ClientKV
	Code   *CodeKV
}

func NewDatabase() *Database {
	return &Database{
		Client: NewClientKV(),
		Code:   NewCodeKV(),
	}
}

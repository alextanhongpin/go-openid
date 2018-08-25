package main

import (
	openid "github.com/alextanhongpin/go-openid"
)

type Database struct {
	Client *ClientKV
	Code   *CodeKV
}

func NewDatabase() *Database {
	client := &ClientKV{
		db: make(map[string]*openid.Client),
	}
	code := &CodeKV{
		db: make(map[string]*openid.Code),
	}
	return &Database{
		Client: client,
		Code:   code,
	}
}

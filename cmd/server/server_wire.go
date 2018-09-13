package main

import (
	"github.com/alextanhongpin/go-openid/internal/database"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
)

func initEndpoints(key string) *Endpoints {
	mem := database.NewInMem()
	cry := crypto.New(key)
	svc := NewService(mem, cry)
	return NewEndpoints(svc)
}

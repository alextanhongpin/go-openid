package main

import (
	"log"

	"github.com/alextanhongpin/go-openid/internal/client"
	"github.com/alextanhongpin/go-openid/internal/core"
	"github.com/alextanhongpin/go-openid/internal/user"
	"github.com/alextanhongpin/go-openid/service"
)

type serviceImpl struct {
	client service.Client
	core   service.Core
	user   service.User
}

// NewService returns a new service.
func NewService() *serviceImpl {
	u := user.NewService()
	c, err := client.NewService()
	if err != nil {
		log.Fatal(err)
	}
	cr := core.NewService()
	if err != nil {
		log.Fatal(err)
	}
	return &serviceImpl{
		client: c,
		user:   u,
		core:   cr,
	}
}

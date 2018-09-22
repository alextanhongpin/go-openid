package main

import (
	"log"

	"github.com/alextanhongpin/go-openid/internal/client"
	"github.com/alextanhongpin/go-openid/internal/user"
	"github.com/alextanhongpin/go-openid/pkg/service"
)

type serviceImpl struct {
	client service.Client
	user   service.User
}

// NewService returns a new service.
func NewService() *serviceImpl {
	u := user.NewService()
	c, err := client.NewService()
	if err != nil {
		log.Fatal(err)
	}
	return &serviceImpl{
		client: c,
		user:   u,
	}
}

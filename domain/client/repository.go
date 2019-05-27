package client

import openid "github.com/alextanhongpin/go-openid"

type Repository interface {
	Create(client openid.Client) (string, error)
	WithCredentials(clientID, clientSecret string) (*openid.Client, error)
}

package client

import openid "github.com/alextanhongpin/go-openid"

type Repository interface {
	Create(client openid.Client) (bool, error)

	// WithID(id string) (*openid.Client, error)
	// WithCredentials(clientID, clientSecret string) (*openid.Client, error)
}

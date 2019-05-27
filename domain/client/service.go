package client

import openid "github.com/alextanhongpin/go-openid"

type Service interface {
	Client() openid.Client
	ValidateCredentials(clientID, clientSecret string) error
	ProvideCredentials(*openid.Client) error
}

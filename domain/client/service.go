package client

import openid "github.com/alextanhongpin/go-openid"

type Service interface {
	ClientSecret() string
	ClientID() string
	AccessToken(clientID string) (string, error)
	Client() openid.Client
	ValidateCredentials(clientID, clientSecret string) error
}

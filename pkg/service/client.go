package service

import "github.com/alextanhongpin/go-openid"

// Client represents the client service.
type Client interface {
	Register(client *oidc.Client) (*oidc.Client, error)
}

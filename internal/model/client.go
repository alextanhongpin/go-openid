package model

import "github.com/alextanhongpin/go-openid"

// Client represents the client model.
type Client interface {
	New(client *oidc.Client) *oidc.Client
	Save(client *oidc.Client) (*oidc.Client, error)
}

package service

import (
	"errors"

	"github.com/alextanhongpin/go-openid/domain/client"
	"github.com/alextanhongpin/go-openid/pkg/gostrings"
	"github.com/alextanhongpin/go-openid/pkg/randstr"
	"github.com/rs/xid"
)

type Client struct {
	clients client.Repository
}

func (c *Client) Validate(client openid.Client) error {
	if gostrings.IsEmpty(client.ClientID) {
		return errors.New("client_id is required")
	}
	if gostrings.IsEmpty(client.ClientSecret) {
		return errors.New("client_secret is required")
	}
	cli, err := c.clients.WithCredentials(client.ClientID, client.ClientSecret)
	if err != nil {
		return err
	}
	// if cli.GetRedirectURIs.Contains(client.RedirectURIs[0])
	return nil
}

func (c *Client) ClientSecret() string {
	return randstr.RandomString()
}

func (c *Client) ClientID() string {
	return xid.New()
}

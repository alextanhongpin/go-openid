package service

import (
	"errors"
	"time"

	"github.com/alextanhongpin/go-openid/domain/client"
	"github.com/alextanhongpin/go-openid/pkg/gostrings"
	"github.com/alextanhongpin/go-openid/pkg/randstr"
	"github.com/alextanhongpin/go-openid/pkg/signer"
	"github.com/rs/xid"
)

type Client struct {
	clients client.Repository
	signer  *signer.Signer
}

func NewClient(clients client.Repository, accessTokenDuration, refreshTokenDuration time.Duration, signer *signer.Signer) *Client {
	return &Client{
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
		clients:              clients,
		signer:               signer,
	}
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

func (c *Client) ProvideCredentials(client *openid.Client) {
	now := time.Now().UTC()
	client.ClientID = xid.New().String()
	client.ClientSecret = randstr.RandomString(32)
	client.RegistrationAccessToken = ""
	client.ClientIDIssuedAt = now.Unix()
	client.ClientSecretExpiresAt = 0
	client.RegistrationClientURI = "https://server.example.com/c2id/clients"
	{
		var (
			aud = "https://server.example.com/c2id/clients"
			sub = clientID
			iss = clientID

			iat = now
			day = time.Hour * 24
			exp = iat.Add(7 * day)
			key = []byte("client_token_secret")
		)
		claims := signer.NewStandardClaims(aud, sub, iss, iat.Unix(), exp.Unix())
		accessToken, err := c.signer.Sign(claims)
		client.RegistrationAccessToken = accessToken
		return err
	}
	return nil
}

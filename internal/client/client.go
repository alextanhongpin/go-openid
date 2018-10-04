package client

import (
	"time"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
)

// NewClient returns a new client with the generated client id and secret. It basically clones an existing client, which can be unmarshalled from json, and injects the credentials.
func NewClient(c *openid.Client) (*openid.Client, error) {
	client := c.Clone()

	// Generate client id.
	clientID := crypto.NewXID()

	// Generate client secret.
	clientSecret, err := crypto.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}

	// TODO: Move secret to envvars.
	var (
		iat = time.Now().UTC()
		day = time.Hour * 24
		exp = iat.Add(7 * day)
		aud = "https://server.example.com/c2id/clients"
		iss = clientID
		sub = clientID
		key = []byte("client_token_secret")
	)

	// Generate access token.
	claims := crypto.NewStandardClaims(aud, sub, iss, iat.Unix(), exp.Unix())
	accessToken, err := crypto.NewJWT(key, claims)
	if err != nil {
		return nil, err
	}

	client.ClientID = clientID
	client.ClientSecret = clientSecret
	client.ClientIDIssuedAt = iat.Unix()
	client.ClientSecretExpiresAt = 0 // Never expire.
	client.RegistrationAccessToken = accessToken
	client.RegistrationClientURI = aud

	return client, nil
}

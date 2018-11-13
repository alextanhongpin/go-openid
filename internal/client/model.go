package client

import (
	"time"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
	"github.com/alextanhongpin/go-openid/pkg/schema"
)

type Model struct {
}

// NewModel returns a new client model implementation.
func NewModel() *Model {
	return &Model{}
}

func (m *Model) GenerateClientID() string {
	return crypto.NewXID()
}

func (m *Model) GenerateClientSecret() (string, error) {
	secret, err := crypto.GenerateRandomString(32)
	return secret, err
}

func (m *Model) GenerateRegistrationAccessToken(clientID string) (string, error) {
	var (
		aud = "https://server.example.com/c2id/clients"
		sub = clientID
		iss = clientID

		iat = time.Now().UTC()
		day = time.Hour * 24
		exp = iat.Add(7 * day)
		key = []byte("client_token_secret")
	)
	claims := crypto.NewStandardClaims(aud, sub, iss, iat.Unix(), exp.Unix())
	accessToken, err := crypto.NewJWT(key, claims)
	return accessToken, err
}

func (m *Model) GenerateClientIDIssuedAt() int64 {
	return time.Now().UTC().Unix()
}

func (m *Model) GenerateClientSecretExpiresAt() int64 {
	return 0
}

func (m *Model) GenerateRegistrationClientURI() string {
	return "https://server.example.com/c2id/clients"
}

func (m *Model) ValidateClient(o *openid.Client) error {
	validator, err := schema.NewClientValidator()
	if err != nil {
		return err
	}
	_, err = validator.Validate(o)
	return err
}

func (m *Model) ValidateClientResponse(o *openid.Client) error {
	validator, err := schema.NewClientResponseValidator()
	if err != nil {
		return err
	}
	_, err = validator.Validate(o)
	return err
}

type ClientModifier func(o *openid.Client) error

func ClientID(fn func() string) ClientModifier {
	return func(o *openid.Client) error {
		o.ClientID = fn()
		return nil
	}
}

func ClientSecret(fn func() (string, error)) ClientModifier {
	return func(o *openid.Client) error {
		var err error
		o.ClientSecret, err = fn()
		return err
	}
}

func RegistrationAccessToken(fn func(string) (string, error)) ClientModifier {
	return func(o *openid.Client) error {
		var err error
		o.RegistrationAccessToken, err = fn(o.ClientID)
		return err
	}
}

func ClientIDIssuedAt(fn func() int64) ClientModifier {
	return func(o *openid.Client) error {
		o.ClientIDIssuedAt = fn()
		return nil
	}
}

func ClientSecretExpiresAt(fn func() int64) ClientModifier {
	return func(o *openid.Client) error {
		o.ClientSecretExpiresAt = fn()
		return nil
	}
}
func RegistrationClientURI(fn func() string) ClientModifier {
	return func(o *openid.Client) error {
		o.RegistrationClientURI = fn()
		return nil
	}
}

func ClientResponseValidation(fn func(o *openid.Client) error) ClientModifier {
	return func(o *openid.Client) error {
		return fn(o)
	}
}
func ClientValidation(fn func(o *openid.Client) error) ClientModifier {
	return func(o *openid.Client) error {
		return fn(o)
	}
}

func apply(o *openid.Client, modifiers ...ClientModifier) error {
	var err error
	for _, m := range modifiers {
		err = m(o)
		if err != nil {
			return err
		}
	}
	return err
}

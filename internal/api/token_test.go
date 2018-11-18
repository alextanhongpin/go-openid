package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type codeRepository struct {
	code *Code
}

func (c *codeRepository) Create(*Code) error {
	return nil
}
func (c *codeRepository) GetCodeByID(id string) (*Code, error) {
	if c.code.ID != id {
		return nil, errors.New("code does not exist")
	}
	return c.code, nil
}

func TestTokenFlow(t *testing.T) {
	assert := assert.New(t)
	var (
		clientID     = "client_id"
		clientSecret = "client_secret"
		code         = "xyz"
		subject      = "user_1"
		redirectURI  = "https://client.example.com/cb"
		now          = time.Now()
	)
	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextKeyClientID, clientID)
	ctx = context.WithValue(ctx, ContextKeyClientSecret, clientSecret)
	ctx = context.WithValue(ctx, ContextKeySubject, subject)
	ctx = context.WithValue(ctx, ContextKeyTimestamp, now)

	clientRepo := new(clientRepository)
	clientRepo.client = NewClient()
	clientRepo.client.RedirectURIs = append(clientRepo.client.RedirectURIs, redirectURI)

	codeRepo := new(codeRepository)
	codeRepo.code = NewCode(code, 10*time.Second)

	req := &TokenRequest{
		GrantType:   "authorization_code",
		Code:        code,
		RedirectURI: redirectURI,
	}
	signer := NewNopSigner()
	res, err := Token(ctx, clientRepo, codeRepo, signer, req)
	if err != nil {
		t.Fatalf("want error nil, got %v", err)
	}
	accessTokenClaims, err := signer.Parse(res.AccessToken)
	assert.Nil(err)
	assert.Equal(subject, accessTokenClaims.Subject)
	assert.Equal(now.Add(2*time.Hour).Unix(), accessTokenClaims.ExpiresAt)
}

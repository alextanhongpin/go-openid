package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

type codeRepository struct {
	code *Code
}

func (c *codeRepository) GetCodeByID(id string) (*Code, error) {
	if c.code.ID != id {
		return nil, errors.New("code does not exist")
	}
	return c.code, nil
}

type tokenResponseBuilderMock struct{}

func (t *tokenResponseBuilderMock) Build(*TokenRequest) (*TokenResponse, error) {
	return NewTokenResponse("access_token", "refresh_token", "id_token", 3600), nil
}

func TestTokenFlow(t *testing.T) {
	var (
		clientID     = "client_id"
		clientSecret = "client_secret"
		code         = "xyz"
	)
	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextKeyClientID, clientID)
	ctx = context.WithValue(ctx, ContextKeyClientSecret, clientSecret)

	clientRepo := new(clientRepository)
	clientRepo.client = NewClient()
	clientRepo.client.RedirectURIs = append(clientRepo.client.RedirectURIs, "https://client.example.com/cb")

	codeRepo := new(codeRepository)
	codeRepo.code = NewCode(code, 10*time.Second)

	req := &TokenRequest{
		GrantType:   "authorization_code",
		Code:        code,
		RedirectURI: "https://client.example.com/cb",
	}
	// responseBuilder := NewTokenResponseBuilder(3600, []byte("secret"))
	responseBuilder := new(tokenResponseBuilderMock)
	res, err := Token(ctx, clientRepo, codeRepo, req, responseBuilder)
	if err != nil {
		t.Fatalf("want error nil, got %v", err)
	}
	if accessToken := res.AccessToken; accessToken != "access_token" {
		t.Fatalf("want %v, got %v", "access_token", accessToken)
	}
	if refreshToken := res.RefreshToken; refreshToken != "refresh_token" {
		t.Fatalf("want %v, got %v", "refresh_token", refreshToken)
	}
	if idToken := res.IDToken; idToken != "id_token" {
		t.Fatalf("want %v, got %v", "id_token", idToken)
	}

}

package main

import (
	"errors"
	"testing"
	"time"
)

type clientRepository struct {
	client *Client
}

func (c *clientRepository) GetClientByClientID(clientID string) (*Client, error) {
	if c.client.ClientID != clientID {
		return nil, errors.New("client_id not found")
	}
	return c.client, nil
}

func (c *clientRepository) GetClientByCredentials(clientID, clientSecret string) (*Client, error) {
	return c.client, nil
}

func TestAuthenticate(t *testing.T) {
	var (
		scope       = "openid profile"
		clientID    = "1"
		redirectURI = "https://client.example.com/cb"
		state       = "xyz"
		code        = "c0d3"
	)
	// Prepare request.
	request := NewAuthenticateRequest(scope, clientID, redirectURI)
	request.State = state

	// Prepare repository.
	clientRepo := new(clientRepository)
	clientRepo.client = NewClient()
	clientRepo.client.ClientID = clientID
	clientRepo.client.RedirectURIs = append(clientRepo.client.RedirectURIs, redirectURI)

	codeRepo := new(codeRepository)
	codeFactory := func() *Code {
		return NewCode(code, 10*time.Minute)
	}
	// Call service.
	response, err := Authenticate(clientRepo, codeRepo, codeFactory, request)
	if err != nil {
		t.Fatalf("want error nil, got %v", err)
	}
	if got := response.Code; got != code {
		t.Fatalf("want %v, got %v", code, got)
	}
	if got := response.State; got != state {
		t.Fatalf("want %v, got %v", state, got)
	}
}

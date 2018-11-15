package main

import (
	"testing"
)

type clientRepository struct {
	client *Client
}

func (c *clientRepository) GetClientByClientID(clientID string) (*Client, error) {
	return c.client, nil
}

func (c *clientRepository) GetClientByCredentials(clientID, clientSecret string) (*Client, error) {
	return c.client, nil
}

func TestAuthenticate(t *testing.T) {
	// Prepare request.
	request := NewAuthenticateRequest("openid", "client_id", "redirect_uri")
	request.State = "xyz"

	// Prepare repository.
	clientRepo := new(clientRepository)
	clientRepo.client = NewClient()
	clientRepo.client.RedirectURIs = append(clientRepo.client.RedirectURIs, "redirect_uri")

	// Call service.
	response, err := AuthenticateFlow(clientRepo, request, func() string {
		return "new_code"
	})
	if err != nil {
		t.Fatalf("want error nil, got %v", err)
	}
	if code := response.Code; code != "new_code" {
		t.Fatalf("want %v, got %v", "new_code", code)
	}
	if state := response.State; state != "xyz" {
		t.Fatalf("want %v, got %v", "xyz", state)
	}
}

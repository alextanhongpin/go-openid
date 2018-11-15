package main

import (
	"errors"
	"fmt"
	"strings"
)

type (
	AuthenticateRequest struct {
		Scope        string
		ResponseType string
		ClientID     string
		RedirectURI  string
		State        string
	}

	AuthenticateResponse struct {
		Code  string
		State string
	}
)

func NewAuthenticateRequest(scope, clientID, redirectURI string) *AuthenticateRequest {
	return &AuthenticateRequest{
		Scope:        scope,
		ResponseType: "code",
		ClientID:     clientID,
		RedirectURI:  redirectURI,
	}
}
func NewAuthenticateResponse(code, state string) *AuthenticateResponse {
	return &AuthenticateResponse{
		Code:  code,
		State: state,
	}
}

func Authenticate(repo ClientRepository, req *AuthenticateRequest) error {
	// Pre-Authenticate: Performs validations.
	if err := ValidateAuthenticateRequest(req); err != nil {
		return err
	}
	// Do work
	client, err := ValidateClient(repo, req.ClientID)
	if err != nil {
		return err
	}
	if !URIs(client.RedirectURIs).Contains(req.RedirectURI) {
		return errors.New("redirect_uri is invalid")
	}
	// Post-Authenticate: Builds response.
	return nil
}

func PostAuthenticate(state string, code func() string) *AuthenticateResponse {
	return NewAuthenticateResponse(code(), state)
}

// Full/Partial? Partial will not return any response. By default, it will always be treated as full, and hence does not need the suffix -Full. AuthenticatePartial if no response is required. If there are multiple steps, can add the suffix -Flow.
func AuthenticateFlow(repo ClientRepository, req *AuthenticateRequest, codeGenerator func() string) (*AuthenticateResponse, error) {
	// Combines both pre-authenticate and do work.
	if err := Authenticate(repo, req); err != nil {
		return nil, err
	}
	res := PostAuthenticate(req.State, codeGenerator)
	return res, nil
}

func ValidateAuthenticateRequest(req *AuthenticateRequest) error {
	// Validate required fields.
	fields := []struct {
		label, value string
	}{
		{"scope", req.Scope},
		{"response_type", req.ResponseType},
		{"client_id", req.ClientID},
		{"redirect_uri", req.RedirectURI},
	}
	for _, field := range fields {
		if stringIsEmpty(field.value) {
			return fmt.Errorf(`"%s" is required`, field.label)
		}
	}
	// Validate specific fields.
	if !strings.Contains(req.Scope, "openid") {
		return errors.New(`scope "openid" is required`)
	}
	if !strings.EqualFold(req.ResponseType, "code") {
		return errors.New(`response_type "code" is invalid`)
	}
	return nil
}

func ValidateClient(repo ClientRepository, clientID string) (*Client, error) {
	return repo.GetClientByClientID(clientID)
}

func ValidateClientCredentials(repo ClientRepository, clientID, clientSecret string) (*Client, error) {
	return repo.GetClientByCredentials(clientID, clientSecret)
}

package main

import (
	"errors"
	"fmt"
	"strings"
)

type (
	// AuthenticateRequest represents the authentication request.
	AuthenticateRequest struct {
		Scope        string
		ResponseType string
		ClientID     string
		RedirectURI  string
		State        string
	}

	// AuthenticateResponse represents the successfull authentication
	// response.
	AuthenticateResponse struct {
		Code  string
		State string
	}
)

// NewAuthenticateRequest returns a new AuthenticateRequest with the response
// type code.
func NewAuthenticateRequest(scope, clientID, redirectURI string) *AuthenticateRequest {
	return &AuthenticateRequest{
		Scope:        scope,
		ResponseType: "code",
		ClientID:     clientID,
		RedirectURI:  redirectURI,
	}
}

// NewAuthenticateResponse returns a new AuthenticateResponse with a unique,
// code that is valid for a fixed duration.
func NewAuthenticateResponse(code, state string) *AuthenticateResponse {
	return &AuthenticateResponse{
		Code:  code,
		State: state,
	}
}

// Authenticate performs a validation on the request.
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

// PostAuthenticate returns a response with a code. This is called after
// Authenticate method validates the request successfully.
func PostAuthenticate(state string, code func() string) *AuthenticateResponse {
	return NewAuthenticateResponse(code(), state)
}

// NOTE: Full/Partial? Partial will not return any response. By default, it
// will always be treated as full, and hence does not need the suffix -Full.
// AuthenticatePartial if no response is required. If there are multiple steps,
// can add the suffix -Flow.

// AuthenticateFlow represents the authentication flow.
func AuthenticateFlow(
	repo ClientRepository,
	req *AuthenticateRequest,
	codeGenerator func() string,
) (*AuthenticateResponse, error) {
	// Combines both pre-authenticate and do work.
	if err := Authenticate(repo, req); err != nil {
		return nil, err
	}
	res := PostAuthenticate(req.State, codeGenerator)
	return res, nil
}

// ValidateAuthenticateRequest checks for the required fields and the values
// set.
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

// ValidateClient validates the client by checking if the given clientID exists
// in the repository.
func ValidateClient(repo ClientRepository, clientID string) (*Client, error) {
	return repo.GetClientByClientID(clientID)
}

// ValidateClientCredentials checks if the given client credentials matches a
// client in the repository.
func ValidateClientCredentials(
	repo ClientRepository,
	clientID,
	clientSecret string,
) (*Client, error) {
	return repo.GetClientByCredentials(clientID, clientSecret)
}

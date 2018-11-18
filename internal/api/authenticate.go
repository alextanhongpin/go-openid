package main

import (
	"errors"
	"strings"
)

type (
	// AuthenticateRequest represents the authentication request.
	AuthenticateRequest struct {
		AcrValues    string `json:"acr_values,omitempty"`
		ClientID     string `json:"client_id,omitempty"`
		Display      string `json:"display,omitempty"`
		IDTokenHint  string `json:"id_token_hint,omitempty"`
		LoginHint    string `json:"login_hint,omitempty"`
		MaxAge       int64  `json:"max_age,omitempty"`
		Nonce        string `json:"nonce,omitempty"`
		Prompt       string `json:"prompt,omitempty"`
		RedirectURI  string `json:"redirect_uri,omitempty"`
		ResponseMode string `json:"response_mode,omitempty"`
		ResponseType string `json:"response_type,omitempty"`
		Scope        string `json:"scope,omitempty"`
		State        string `json:"state,omitempty"`
		UILocales    string `json:"ui_locales,omitempty"`
	}

	// AuthenticateResponse represents the successfull authentication
	// response.
	AuthenticateResponse struct {
		Code  string `json:"code,omitempty"`
		State string `json:"state,omitempty"`
	}
)

type Scope string

func (s Scope) Has(scope string) bool {
	return strings.Contains(string(s), scope)
}

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

// NOTE: Full/Partial? Partial will not return any response. By default, it
// will always be treated as full, and hence does not need the suffix -Full.
// AuthenticatePartial if no response is required. If there are multiple steps,
// can add the suffix -Flow.

// AuthenticateFlow represents the authentication flow.
func Authenticate(
	repo ClientRepository,
	codeRepo CodeRepository,
	codeFactory CodeFactory,
	req *AuthenticateRequest,
) (*AuthenticateResponse, error) {
	// Pre-Authenticate: Performs validations.
	if err := ValidateAuthenticateRequest(req); err != nil {
		return nil, err
	}
	// Can it be simplified to get client by clientID and redirectURI?
	client, err := ValidateClient(repo, req.ClientID)
	if err != nil {
		return nil, err
	}
	if !URIs(client.RedirectURIs).Contains(req.RedirectURI) {
		return nil, errors.New("redirect_uri is invalid")
	}
	// Create a code and store it.
	code := codeFactory()
	if err := CreateCode(codeRepo, code); err != nil {
		return nil, err
	}
	res := NewAuthenticateResponse(code.ID, req.State)
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
			return MakeErrRequired(field.label)
		}
	}
	// Validate specific fields.
	// Hide implementation details with structs.
	if !Scope(req.Scope).Has("openid") {
		return MakeErrRequired("openid")
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

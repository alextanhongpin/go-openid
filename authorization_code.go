package oidc

import "strings"

// AuthorizationRequest represents the request payload for authorization.
type AuthorizationRequest struct {
	ResponseType string `json:"response_type,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	RedirectURI  string `json:"redirect_uri,omitempty"`
	Scope        string `json:"scope,omitempty"`
	State        string `json:"state,omitempty"`
}

// Validate performs validation on required fields.
func (r *AuthorizationRequest) Validate() error {
	// Required fields
	if r.ResponseType != "code" {
		return ErrUnsupportedResponseType
	}
	if strings.TrimSpace(r.ClientID) == "" {
		return ErrUnauthorizedClient
	}
	// Optional fields
	if strings.TrimSpace(r.RedirectURI) == "" {
		return ErrInvalidRequest
	}
	if strings.TrimSpace(r.Scope) == "" {
		return ErrInvalidScope
	}
	return nil
}

// AuthorizationResponse represents the authorization response body.
type AuthorizationResponse struct {
	Code  string `json:"code,omitempty"`
	State string `json:"state,omitempty"`
}

// AuthorizationError represents the struct for the error.
type AuthorizationError struct {
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
	ErrorURI         string `json:"error_uri,omitempty"`
	State            string `json:"state,omitempty"`
}

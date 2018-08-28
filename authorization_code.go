package oidc

import (
	"net/url"
)

// AuthorizationRequest represents the request payload for authorization.
type AuthorizationRequest struct {
	ResponseType string `json:"response_type"`
	ClientID     string `json:"client_id"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
	State        string `json:"state"`
}

// DecodeAuthorizationRequest takes in a url with the query string parameters
// and converts it into a struct.
func DecodeAuthorizationRequest(u url.Values) *AuthorizationRequest {
	return &AuthorizationRequest{
		ResponseType: u.Get("response_type"),
		ClientID:     u.Get("client_id"),
		RedirectURI:  u.Get("redirect_uri"),
		Scope:        u.Get("scope"),
		State:        u.Get("state"),
	}
}

// EncodeAuthorizationRequest converts the struct into url with query string.
func EncodeAuthorizationRequest(r *AuthorizationRequest) (url.Values, error) {
	u, err := url.Parse("")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("response_type", r.ResponseType)
	q.Add("client_id", r.ClientID)
	q.Add("redirect_uri", r.RedirectURI)
	q.Add("scope", r.Scope)
	q.Add("state", r.State)

	return q, nil
}

// Validate performs validation on required fields.
func (r *AuthorizationRequest) Validate() error {
	// Required fields
	if r.ResponseType != "code" {
		return ErrUnsupportedResponseType
	}
	if r.ClientID == "" {
		return ErrUnauthorizedClient
	}
	// Optional fields
	if r.RedirectURI == "" {
		return ErrInvalidRequest
	}
	if r.Scope == "" {
		return ErrInvalidScope
	}
	// if r.State == "" { }
	return nil
}

type AuthorizationResponse struct {
	Code  string `json:"code,omitempty"`
	State string `json:"state,omitempty"`
}

func DecodeAuthorizationResponse(u url.Values) *AuthorizationResponse {
	return &AuthorizationResponse{
		Code:  u.Get("code"),
		State: u.Get("state"),
	}
}
func EncodeAuthorizationResponse(r *AuthorizationResponse, targetURL string) (*url.URL, error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("code", r.Code)
	q.Add("state", r.State)
	u.RawQuery = q.Encode()
	return u, nil
}

// AuthorizationError represents the struct for the error.
type AuthorizationError struct {
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
	ErrorURI         string `json:"error_uri,omitempty"`
	State            string `json:"state,omitempty"`
}

// EncodeAuthorizationError takes a struct and url and embed the struct as query string parameters to the url.
func EncodeAuthorizationError(r *AuthorizationError, targetURL string) (*url.URL, error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("error", r.Error)
	q.Add("error_description", r.ErrorDescription)
	q.Add("error_uri", r.ErrorURI)
	q.Add("state", r.State)
	u.RawQuery = q.Encode()
	return u, nil
}

// DecodeAuthorizationRequest takes the query string and returns a struct.
func DecodeAuthorizationError(u url.Values) *AuthorizationError {
	return &AuthorizationError{
		Error:            u.Get("error"),
		ErrorDescription: u.Get("error_description"),
		ErrorURI:         u.Get("error_uri"),
		State:            u.Get("state"),
	}
}

package oidc

import (
	"errors"
)

const (
	// Bearer represents the bearer type token.
	Bearer = "Bearer"
	Basic  = "Basic"
)

//go:generate gencodec -type AccessTokenRequest -out gen_token_json.go

// AccessTokenRequest represents the access token request payload.
type AccessTokenRequest struct {
	GrantType   string `json:"grant_type,omitempty"`
	Code        string `json:"code,omitempty"`
	RedirectURI string `json:"redirect_uri,omitempty"`
	ClientID    string `json:"client_id,omitempty"`
}

// Validate performs an initial validation on the required field.
func (r *AccessTokenRequest) Validate() error {
	// if r.GrantType != "authorization_code" {
	//         return InvalidRequest.JSON()
	// }
	// if strings.TrimSpace(r.Code) == "" {
	//         return AccessDenied.JSON()
	// }
	// if strings.TrimSpace(r.RedirectURI) == "" {
	//         return InvalidRequest.JSON()
	// }
	// if !govalidator.IsURL(r.RedirectURI) {
	//         return InvalidRedirectURI.JSON()
	// }
	// if strings.TrimSpace(r.ClientID) == "" {
	//         return AccessDenied.JSON()
	// }
	return nil
}

// AccessTokenResponse represents the response payload.
type AccessTokenResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
}

// RefreshTokenRequest represents the refresh token request.
type RefreshTokenRequest struct {
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	GrantType    string `json:"grant_type,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// Validate checks for required fields.
func (r *RefreshTokenRequest) Validate() error {
	if r.GrantType != "refresh_token" {
		return errors.New("invalid_grant_type")
	}
	if r.RefreshToken == "" {
		return errors.New("invalid_request")
	}
	if r.Scope == "" {
		// TODO: Handle validation for scope.
	}
	return nil
}

// RefreshTokenResponse returns the access token.
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	IDTokens    string `json:"id_token,omitempty"`
}

package oidc

import "strings"

const (
	Bearer = "Bearer"
	Basic  = "Basic"
)

// AccessTokenRequest represents the access token request payload.
type AccessTokenRequest struct {
	GrantType   string `json:"grant_type,omitempty"`
	Code        string `json:"code,omitempty"`
	RedirectURI string `json:"redirect_uri,omitempty"`
	ClientID    string `json:"client_id,omitempty"`
}

// Validate performs an initial validation on the required field.
func (r *AccessTokenRequest) Validate() error {
	if r.GrantType != "authorization_code" {
		return InvalidRequest.JSON()
	}
	if strings.TrimSpace(r.Code) == "" {
		return AccessDenied.JSON()
	}
	if strings.TrimSpace(r.RedirectURI) == "" {
		return InvalidRequest.JSON()
	}
	if strings.TrimSpace(r.ClientID) == "" {
		return AccessDenied.JSON()
	}
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

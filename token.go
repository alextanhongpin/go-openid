package oidc

type AccessTokenRequest struct {
	GrantType   string `json:"grant_type,omitempty"`
	Code        string `json:"code,omitempty"`
	RedirectURI string `json:"redirect_uri,omitempty"`
	ClientID    string `json:"client_id,omitempty"`
}

func (r *AccessTokenRequest) Validate() error {
	if r.GrantType != "authorization_code" {
		return ErrInvalidRequest
	}
	// Check required field
	if r.Code == "" {
		return ErrForbidden
	}

	if r.RedirectURI == "" {
		return ErrForbidden
	}
	if r.ClientID == "" {
		return ErrForbidden
	}
	return nil
}

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

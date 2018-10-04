package openid

import (
	"net/url"

	"github.com/alextanhongpin/go-openid/pkg/querystring"
)

//go:generate gencodec -type AuthenticationRequest -out gen_authn.go

// TODO: Document the helpers difference
// verify - perform comparison of value; returns true or false
// validate - check for required fields or incorrect verification; returns error
// https://www.easterbrook.ca/steve/2010/11/the-difference-between-verification-and-validation/

// AuthenticationRequest is an OAuth 2.0 Authorization Request that requests
// that the End User be authenticated by the Authorization Server.
type AuthenticationRequest struct {
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

// GetPrompt returns the prompt.
func (a *AuthenticationRequest) GetPrompt() Prompt {
	return parsePrompt(a.Prompt)
}

// GetResponseType returns the response type as bitwise int.
func (a *AuthenticationRequest) GetResponseType() ResponseType {
	return parseResponseType(a.ResponseType)
}

// GetScope returns the scope as bitwise int.
func (a *AuthenticationRequest) GetScope() Scope {
	return parseScope(a.Scope)
}

// FromQueryString decodes an authentication request from the given querystring.
func (a *AuthenticationRequest) FromQueryString(u url.Values) error {
	return querystring.Decode(u, a)
}

// AuthenticationResponse is an OAuth 2.0 Authorization Response message
// returned from the OP's Authorization Endpoint in response to the
// Authorization Request message sent by the RP.
type AuthenticationResponse struct {
	Code  string `json:"code,omitempty"`
	State string `json:"state,omitempty"`
}

// ToQueryString converts the response struct into url.Values.
func (a *AuthenticationResponse) ToQueryString() url.Values {
	return querystring.Encode(url.Values{}, a)
}

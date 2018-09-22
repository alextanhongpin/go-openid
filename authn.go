package oidc

import (
	"errors"

	"github.com/asaskevich/govalidator"
)

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

func (a *AuthenticationRequest) GetPrompt() Prompt {
	return parsePrompt(a.Prompt)
}

// Validate performs validation on the required fields except the validation of
// the client, which requires a call to another domain service.
func (a *AuthenticationRequest) Validate() error {
	if err := a.ValidateScope(); err != nil {
		return err
	}
	if err := a.ValidateResponseType(); err != nil {
		return err
	}
	return nil
}

func (a *AuthenticationRequest) ValidateScope() error {
	// REQUIRED. OpenID Connect requests MUST contain the openid scope
	// value. If the openid scope value is not present, the behavior is
	// entirely unspecified. Other scope values MAY be present. Scope
	// values used that are not understood by an implementation SHOULD be
	// ignored. See Sections 5.4 and 11 for additional scope values defined
	// by this specification.
	return validateScope(a.Scope)
}

func (a *AuthenticationRequest) ValidateResponseType() error {
	// REQUIRED. OAuth 2.0 Response Type value that determines the
	// authorization processing flow to be used, including what parameters
	// are returned from the endpoints used. When using the Authorization
	// Code Flow, this value is code.
	return validateResponseType(a.ResponseType)
}

type ClientValidator func(clientID string) error

func (a *AuthenticationRequest) VerifyClientID(validator ClientValidator) error {
	// REQUIRED. OAuth 2.0 Client Identifier valid at the Authorization
	// Server.
	return validator(a.ClientID)
}

func (a *AuthenticationRequest) VerifyRedirectURI(uris RedirectURIs) error {
	// REQUIRED. Redirection URI to which the response will be sent. This
	// URI MUST exactly match one of the Redirection URI values for the
	// Client pre-registered at the OpenID Provider, with the matching
	// performed as described in Section 6.2.1 of [RFC3986] (Simple String
	// Comparison). When using this flow, the Redirection URI SHOULD use
	// the https scheme; however, it MAY use the http scheme, provided that
	// the Client Type is confidential, as defined in Section 2.1 of OAuth
	// 2.0, and provided the OP allows the use of http Redirection URIs in
	// this case. The Redirection URI MAY use an alternate scheme, such as
	// one that is intended to identify a callback into a native
	// application.
	return validateRedirectURI(a.RedirectURI, uris)
}

func (a *AuthenticationRequest) VerifyState(state string) error {
	// RECOMMENDED. Opaque value used to maintain state between the request
	// and the callback. Typically, Cross-Site Request Forgery (CSRF, XSRF)
	// mitigation is done by cryptographically binding the value of this
	// parameter with a browser cookie.
	if a.State != state {
		return errors.New("invalid state")
	}
	return nil
}

// AuthenticationResponse is an OAuth 2.0 Authorization Response message
// returned from the OP's Authorization Endpoint in response to the
// Authorization Request message sent by the RP.
type AuthenticationResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
	IDToken     string `json:"id_token,omitempty"`
	State       string `json:"state,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
}

// -- helpers

func validateScope(scope string) error {
	if scope == "" {
		return errors.New("scope required")
	}
	if !parseScope(scope).Has(ScopeOpenID) {
		return errors.New("invalid scope")
	}
	return nil
}

func validateResponseType(responseType string) error {
	if responseType == "" {
		return errors.New("response_type required")
	}

	var (
		a = ResponseTypeCode
		b = ResponseTypeIDToken
		c = ResponseTypeToken
	)

	parsed := parseResponseType(responseType)
	if !parsed.Has(a | b | c) {
		return errors.New("invalid response_type")
	}
	return nil
}

func validateRedirectURI(redirectURI string, list RedirectURIs) error {
	if !govalidator.IsURL(redirectURI) {
		return errors.New("redirect_uri invalid")
	}
	if !list.Contains(redirectURI) {
		return errors.New("redirect_uri does not match")
	}
	return nil
}

func verifyDisplay(display string) bool {
	_, ok := displaymap[display]
	return ok
}

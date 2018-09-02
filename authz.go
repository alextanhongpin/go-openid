package oidc

import "strings"

type ErrorCode int

const (
	InvalidRequest ErrorCode = iota
	UnauthorizedClient
	AccessDenied
	UnsupportedResponseType
	InvalidScope
	ServerError
	TemporarilyUnavailable
)

var codeErrors = map[ErrorCode]string{
	InvalidRequest:          "the request is missing a required parameter, includes an invalid parameter value more than once, or is otherwise malformed",
	UnauthorizedClient:      "the client is not authorized to request an authorization code using this method",
	AccessDenied:            "the resource owner or authorization server denied the request",
	UnsupportedResponseType: "the authorization server does not support obtaining an authorization code using this method",
	InvalidScope:            "the requested scope is invalid, unknown or malformed",
	ServerError:             "the authorization server encoutered an unexpected condition that prevented it from fulfilling the request",
	TemporarilyUnavailable:  "the authorization server is unable to handle the request due to a temporary overloading or maintenance of the server",
}

func (e ErrorCode) String() string {
	return [...]string{
		"invalid_request",
		"unauthorized_client",
		"access_denied",
		"unsupported_response_type",
		"invalid_scope",
		"server_error",
		"temporarily_unavailable",
	}[e]
}

func (e ErrorCode) JSON() *ErrorJSON {
	return &ErrorJSON{
		Code:        e.String(),
		Description: codeErrors[e],
		URI:         "",
		State:       "",
	}
}

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
		return UnsupportedResponseType.JSON()
	}
	if strings.TrimSpace(r.ClientID) == "" {
		return AccessDenied.JSON()
	}
	// Optional fields
	if strings.TrimSpace(r.RedirectURI) == "" {
		return InvalidRequest.JSON()
	}
	if strings.TrimSpace(r.Scope) == "" {
		return InvalidScope.JSON()
	}
	return nil
}

// AuthorizationResponse represents the authorization response body.
type AuthorizationResponse struct {
	Code  string `json:"code,omitempty"`
	State string `json:"state,omitempty"`
}

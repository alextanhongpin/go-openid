package oidc

import (
	"strings"

	"github.com/asaskevich/govalidator"
)

// ErrorCode represents the ErrorCode.
type ErrorCode int

const (
	AccessDenied ErrorCode = iota
	InvalidRequest
	InvalidScope
	ServerError
	TemporarilyUnavailable
	UnauthorizedClient
	UnsupportedResponseType
)

var errorCodeDescriptions = map[ErrorCode]string{
	AccessDenied:            "the resource owner or authorization server denied the request",
	InvalidRequest:          "the request is missing a required parameter, includes an invalid parameter value more than once, or is otherwise malformed",
	InvalidScope:            "the requested scope is invalid, unknown or malformed",
	ServerError:             "the authorization server encoutered an unexpected condition that prevented it from fulfilling the request",
	TemporarilyUnavailable:  "the authorization server is unable to handle the request due to a temporary overloading or maintenance of the server",
	UnauthorizedClient:      "the client is not authorized to request an authorization code using this method",
	UnsupportedResponseType: "the authorization server does not support obtaining an authorization code using this method",
}

var errorCodes = map[ErrorCode]string{
	AccessDenied:            "access_denied",
	InvalidRequest:          "invalid_request",
	InvalidScope:            "invalid_scope",
	ServerError:             "server_error",
	TemporarilyUnavailable:  "temporarily_unavailable",
	UnauthorizedClient:      "unauthorized_client",
	UnsupportedResponseType: "unsupported_response_type",
}

// String fulfills the stringer method.
func (e ErrorCode) String() string {
	return errorCodes[e]
}

// Description return the general description based on the error code.
func (e ErrorCode) Description() string {
	return errorCodeDescriptions[e]
}

// JSON returns the error as json struct.
func (e ErrorCode) JSON() *ErrorJSON {
	return &ErrorJSON{
		Code:        e.String(),
		Description: e.Description(),
		URI:         "",
		State:       "",
	}
}

// AuthorizationRequest represents the request payload for authorization.
type AuthorizationRequest struct {
	ClientID     string `json:"client_id,omitempty"`
	RedirectURI  string `json:"redirect_uri,omitempty"`
	ResponseType string `json:"response_type,omitempty"`
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
		return InvalidRedirectURI.JSON()
	}
	if !govalidator.IsURL(r.RedirectURI) {
		return InvalidRedirectURI.JSON()
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

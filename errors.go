package openid

import (
	"fmt"
)

const (
	AccessDenied            = "access_denied"
	InvalidRequest          = "invalid_request"
	InvalidScope            = "invalid_scope"
	ServerError             = "server_error"
	TemporarilyUnavailable  = "temporarily_unavailable"
	UnauthorizedClient      = "unauthorized_client"
	UnsupportedResponseType = "unsupported_response_type"
)

// Authorization errors
var (
	ErrAccessDenied            = NewError("access_denied")
	ErrInvalidRequest          = NewError("invalid_request")
	ErrInvalidScope            = NewError("invalid_scope")
	ErrServerError             = NewError("server_error")
	ErrTemporarilyUnavailable  = NewError("temporarily_unavailable")
	ErrUnauthorizedClient      = NewError("unauthorized_client")
	ErrUnsupportedResponseType = NewError("unsupported_response_type")
)

var errorCodeDescriptions = map[string]string{
	AccessDenied:            "the resource owner or authorization server denied the request",
	InvalidRequest:          "the request is missing a required parameter, includes an invalid parameter value more than once, or is otherwise malformed",
	InvalidScope:            "the requested scope is invalid, unknown or malformed",
	ServerError:             "the authorization server encoutered an unexpected condition that prevented it from fulfilling the request",
	TemporarilyUnavailable:  "the authorization server is unable to handle the request due to a temporary overloading or maintenance of the server",
	UnauthorizedClient:      "the client is not authorized to request an authorization code using this method",
	UnsupportedResponseType: "the authorization server does not support obtaining an authorization code using this method",
}

// ErrorText return the general description based on the error code.
func ErrorText(code string) string {
	return errorCodeDescriptions[code]
}

// Client errors.
var (
	// ErrInvalidClientMetadata occurs when the value of one of the client metadata fields is invalid and the server has rejected this request.
	ErrInvalidClientMetadata = NewError("invalid_client_metadata")

	// ErrInvalidRedirectURI occurs when the value of one or more redirect uris is invalid.
	ErrInvalidRedirectURI = NewError("invalid_redirect_uri")
)

// NewError returns a new custom error.
func NewError(code string) *ErrorJSON {
	desc := errorCodeDescriptions[code]
	return &ErrorJSON{Code: code, Description: desc}
}

// ErrorJSON represents the json error.
type ErrorJSON struct {
	Code        string `json:"error,omitempty"`
	Description string `json:"error_description,omitempty"`
	State       string `json:"state,omitempty"`
	URI         string `json:"error_uri,omitempty"`
}

// Error fulfils the error interface methods.
func (e *ErrorJSON) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Description)
}

// SetURI sets the uri of the error.
func (e *ErrorJSON) SetURI(s string) {
	e.URI = s
}

// SetDescription sets the description of the error.
func (e *ErrorJSON) SetDescription(s string) {
	e.Description = s
}

// SetState sets the state of the error.
func (e *ErrorJSON) SetState(s string) {
	e.State = s
}

func (e *ErrorJSON) WithDescription(s string) error {
	e.SetDescription(s)
	return e
}

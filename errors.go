package oidc

import (
	"fmt"
)

// Authentication Error Response
var (
	ErrAccountSelectionRequired = NewError("account_selection_required")
	ErrConsentRequired          = NewError("consent_required")
	ErrInteractionRequired      = NewError("interaction_required")
	ErrInvalidRequestObject     = NewError("invalid_request_object")
	ErrInvalidRequestURI        = NewError("invalid_request_uri")
	ErrLoginRequired            = NewError("login_required")
	ErrRegistrationNotSupported = NewError("registration_not_supported")
	ErrRequestNotSupported      = NewError("request_not_supported")
	ErrRequestURINotSupported   = NewError("request_uri_not_supported")
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

// Client errors.
var (
	// ErrInvalidClientMetadata occurs when the value of one of the client metadata fields is invalid and the server has rejected this request.
	ErrInvalidClientMetadata = NewError("invalid_client_metadata")

	// ErrInvalidRedirectURI occurs when the value of one or more redirect uris is invalid.
	ErrInvalidRedirectURI = NewError("invalid_redirect_uri")
)

// NewError returns a new custom error.
func NewError(code string) *ErrorJSON {
	return &ErrorJSON{Code: code}
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

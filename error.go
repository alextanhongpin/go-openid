package oidc

import (
	"errors"
	"fmt"
)

// Authentication Error Response
var (
	ErrAccountSelectionRequired = errors.New("account selection required")
	ErrConsentRequired          = errors.New("consent required")
	ErrInteractionRequired      = errors.New("interaction required")
	ErrInvalidRequestObject     = errors.New("invalid request object")
	ErrInvalidRequestURI        = errors.New("invalid request uri")
	ErrLoginRequired            = errors.New("login required")
	ErrRegistrationNotSupported = errors.New("registration not supported")
	ErrRequestNotSupported      = errors.New("request not supported")
	ErrRequestURINotSupported   = errors.New("request uri not supported")
)

// Authorization errors
var (
	ErrAccessDenied            = errors.New("access denied")
	ErrInvalidRequest          = errors.New("invalid request")
	ErrInvalidScope            = errors.New("invalid scope")
	ErrServerError             = errors.New("server error")
	ErrTemporarilyUnavailable  = errors.New("temporarily unavailable")
	ErrUnauthorizedClient      = errors.New("unauthorized client")
	ErrUnsupportedResponseType = errors.New("unsupported response type")
)

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

// SetState set the state of the error.
func (e *ErrorJSON) SetState(s string) {
	e.State = s
}

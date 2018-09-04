package oidc

import (
	"errors"
	"fmt"
)


// Authentication Error Response
var (
	ErrInteractionRequired      = errors.New("interaction required")
	ErrLoginRequired            = errors.New("login required")
	ErrAccountSelectionRequired = errors.New("account selection required")
	ErrConsentRequired          = errors.New("consent required")
	ErrInvalidRequestURI        = errors.New("invalid request uri")
	ErrInvalidRequestObject     = errors.New("invalid request object")
	ErrRequestNotSupported      = errors.New("request not supported")
	ErrRequestURINotSupported   = errors.New("request uri not supported")
	ErrRegistrationNotSupported = errors.New("registration not supported")
)

// Authorization errors
var (
	ErrInvalidRequest          = errors.New("invalid request")
	ErrUnauthorizedClient      = errors.New("unauthorized client")
	ErrAccessDenied            = errors.New("access denied")
	ErrUnsupportedResponseType = errors.New("unsupported response type")
	ErrInvalidScope            = errors.New("invalid scope")
	ErrServerError             = errors.New("server error")
	ErrTemporarilyUnavailable  = errors.New("temporarily unavailable")
)

// ErrorJSON represents the json error.
type ErrorJSON struct {
	Code        string `json:"error,omitempty"`
	Description string `json:"error_description,omitempty"`
	URI         string `json:"error_uri,omitempty"`
	State       string `json:"state,omitempty"`
}

func (e *ErrorJSON) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Description)
}

func (e *ErrorJSON) SetState(s string) {
	e.State = s
}

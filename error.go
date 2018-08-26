package openid

import "errors"

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

var ErrForbidden = errors.New("forbidden request")
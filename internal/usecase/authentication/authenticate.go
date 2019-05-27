package authentication

import (
	"context"
	"errors"

	"github.com/alextanhongpin/go-openid/domain/client"
	"github.com/alextanhongpin/go-openid/domain/code"
	"github.com/alextanhongpin/go-openid/usecase"
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

type UseCase struct {
	clients     client.Repository
	codeService code.Service
}

func (u *UseCase) Authenticate(ctx context.Context, req usecase.AuthenticationRequest) (*usecase.AuthenticationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	client, err := u.clients.WithClientID(req.ClientID)
	if err != nil {
		return nil, err
	}
	if !client.RedirectURIs.Contains(req.RedirectURI) {
		return nil, errors.New("redirect_uri is invalid")
	}
	code, err := u.codeService.Code()
	if err != nil {
		return nil, err
	}
	return &usecae.AuthenticationResponse{
		Code:  code.ID,
		State: req.State,
	}, nil
}

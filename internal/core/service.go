package core

import (
	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/model"
)

type serviceImpl struct {
	model model.Core
}

func (a *serviceImpl) Authenticate(req *oidc.AuthenticationRequest) (*oidc.AuthenticationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	var (
		clientID    = req.ClientID
		redirectURI = req.RedirectURI
	)

	if err := a.model.ValidateClient(clientID, redirectURI); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *serviceImpl) Authorize(req *oidc.AuthenticationRequest) (*oidc.AuthorizationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var (
		clientID    = req.ClientID
		redirectURI = req.RedirectURI
		state       = req.State
	)
	if err := s.model.ValidateClient(clientID, redirectURI); err != nil {
		return nil, err
	}
	// Return the code to be exchanged for a token.
	code := s.model.NewCode()

	return &oidc.AuthorizationResponse{
		Code:  code,
		State: state,
	}, nil

}

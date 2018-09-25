package core

import (
	"context"

	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/model"
)

type serviceImpl struct {
	model model.Core
}

// PreAuthenticate will only check if the authentication request is valid and
// the type of flow it is using.
func (s *serviceImpl) PreAuthenticate(req *oidc.AuthenticationRequest) error {
	if err := s.model.ValidateAuthnRequest(req); err != nil {
		return err
	}

	return s.model.ValidateAuthnClient(req)
}

// Authenticate performs the full authentication and validation of all fields.
func (s *serviceImpl) Authenticate(ctx context.Context, req *oidc.AuthenticationRequest) (*oidc.AuthenticationResponse, error) {
	if err := s.model.ValidateAuthnRequest(req); err != nil {
		return nil, err
	}
	if err := s.model.ValidateAuthnUser(ctx, req); err != nil {
		return nil, err
	}
	if err := s.model.ValidateAuthnClient(req); err != nil {
		return nil, err
	}
	return &oidc.AuthenticationResponse{
		Code:  s.model.NewCode(),
		State: req.State,
	}, nil
}

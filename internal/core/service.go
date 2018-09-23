package core

import (
	"github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/model"
)

type serviceImpl struct {
	model model.Core
}

func (s *serviceImpl) Authenticate(req *oidc.AuthenticationRequest) (*oidc.AuthenticationResponse, error) {
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

	return &oidc.AuthenticationResponse{
		Code:  code,
		State: state,
	}, nil

}

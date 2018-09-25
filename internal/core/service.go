package core

import (
	"context"
	"errors"
	"time"

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

func (s *serviceImpl) Token(ctx context.Context, req *oidc.AccessTokenRequest) (*oidc.AccessTokenResponse, error) {
	auth, ok := oidc.GetAuthContextKey(ctx)
	if !ok {
		return nil, errors.New("missing authorization header")
	}
	client, err := s.model.ValidateClientAuthHeader(auth)
	if err != nil {
		return nil, err
	}
	if ok := client.GetRedirectURIs().Contains(req.RedirectURI); !ok {
		return nil, errors.New("redirect_uri does not match")
	}

	userID, ok := oidc.GetUserIDContextKey(ctx)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	accessToken, err := s.model.ProvideToken(userID, 2*time.Hour)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.model.ProvideToken(userID, 24*7*time.Hour)
	if err != nil {
		return nil, err
	}

	idToken, err := s.model.ProvideIDToken(userID)
	if err != nil {
		return nil, err
	}

	res := oidc.AccessTokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64((2 * time.Hour).Seconds()),
		RefreshToken: refreshToken,
		IDToken:      idToken,
	}
	return &res, nil
}

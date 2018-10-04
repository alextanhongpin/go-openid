




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

// NewService returns a new service.
func NewService(model model.Core) serviceImpl {
	return serviceImpl{model}
}

// SetModel sets the existing model to the given model.
func (s *serviceImpl) SetModel(model model.Core) {
	s.model = model
}

// PreAuthenticate will only check if the authentication request is valid and
// the type of flow it is using.
func (s *serviceImpl) PreAuthenticate(req *openid.AuthenticationRequest) error {
	if req == nil {
		return errors.New("arguments cannot be nil")
	}
	if err := s.model.ValidateAuthnRequest(req); err != nil {
		return err
	}

	return s.model.ValidateAuthnClient(req)
}

// Authenticate performs the full authentication and validation of all fields.
func (s *serviceImpl) Authenticate(ctx context.Context, req *openid.AuthenticationRequest) (*openid.AuthenticationResponse, error) {
	if err := s.model.ValidateAuthnRequest(req); err != nil {
		return nil, err
	}
	if err := s.model.ValidateAuthnUser(ctx, req); err != nil {
		return nil, err
	}
	if err := s.model.ValidateAuthnClient(req); err != nil {
		return nil, err
	}
	return &openid.AuthenticationResponse{
		Code:  s.model.NewCode(),
		State: req.State,
	}, nil
}

func (s *serviceImpl) Token(ctx context.Context, req *openid.AccessTokenRequest) (*openid.AccessTokenResponse, error) {
	auth, ok := openid.GetAuthContextKey(ctx)
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

	userID, ok := openid.GetUserIDContextKey(ctx)
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

	res := openid.AccessTokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64((2 * time.Hour).Seconds()),
		RefreshToken: refreshToken,
		IDToken:      idToken,
	}
	return &res, nil
}

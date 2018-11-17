package main

import (
	"context"
)

type Service struct {
	clientRepo    ClientRepository
	codeRepo      CodeRepository
	codeGenerator func() string
	signer        Signer
}

func New(
	clientRepo ClientRepository,
	codeRepo CodeRepository,
	codeGenerator func() string,
	signer Signer,
) *Service {
	return &Service{
		clientRepo,
		codeRepo,
		codeGenerator,
		signer,
	}
}

// Authenticate performs the authentication flow.
func (s *Service) Authenticate(ctx context.Context, req *AuthenticateRequest) (*AuthenticateResponse, error) {
	return AuthenticateFlow(s.clientRepo, req, s.codeGenerator)
}

func (s *Service) token(ctx context.Context, req *TokenRequest) (*TokenResponse, error) {
	return Token(ctx, s.clientRepo, s.codeRepo, s.signer, req)
}

// Token represents the token flow.
func (s *Service) Token(ctx context.Context, req *TokenRequest) (*TokenResponse, error) {
	res, err := s.token(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

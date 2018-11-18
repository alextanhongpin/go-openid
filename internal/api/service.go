package main

import (
	"context"
)

type Service struct {
	clientRepo  ClientRepository
	codeRepo    CodeRepository
	codeFactory CodeFactory
	signer      Signer
}

func New(
	clientRepo ClientRepository,
	codeRepo CodeRepository,
	codeFactory CodeFactory,
	signer Signer,
) *Service {
	return &Service{
		clientRepo,
		codeRepo,
		codeFactory,
		signer,
	}
}

// Authenticate performs the authentication flow.
func (s *Service) Authenticate(ctx context.Context, req *AuthenticateRequest) (*AuthenticateResponse, error) {
	return Authenticate(
		s.clientRepo,
		s.codeRepo,
		s.codeFactory,
		req,
	)
}

// Token represents the token flow.
func (s *Service) Token(ctx context.Context, req *TokenRequest) (*TokenResponse, error) {
	return Token(
		ctx,
		s.clientRepo,
		s.codeRepo,
		s.signer,
		req,
	)
}

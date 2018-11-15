package main

import "context"

type Service struct {
	clientRepo           ClientRepository
	codeRepo             CodeRepository
	codeGenerator        func() string
	tokenResponseBuilder TokenResponseBuilder
}

func New() *Service {
	return &Service{}
}

func (s *Service) Authenticate(ctx context.Context, req *AuthenticateRequest) (*AuthenticateResponse, error) {
	return AuthenticateFlow(s.clientRepo, req, s.codeGenerator)
}

/*
   u, err := url.Parse(req.RedirectURI)
   if err != nil {}
   q := u.Query()
   q.Set("code", res.Code)
   q.Set("state", res.State)
   u.RawQuery = q.Encode()
   location := u.String()
*/

func (s *Service) token(ctx context.Context, req *TokenRequest) (*TokenResponse, error) {
	return Token(ctx, s.clientRepo, s.codeRepo, req, s.tokenResponseBuilder)
}

func (s *Service) Token(ctx context.Context, req *TokenRequest) (*TokenResponse, error) {
	res, err := s.token(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

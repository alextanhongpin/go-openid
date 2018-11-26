package main

import (
	"context"
	"time"
)

// type ServiceFacade struct {
//         Authenticate AuthenticateService
//         Token        TokenService
// }

// AuthenticateService performs the authentication flow.
type AuthenticateService func(ctx context.Context, req *AuthenticateRequest) (*AuthenticateResponse, error)

// MakeAuthenticateService returns a new AuthenticateService.
func MakeAuthenticateService(repo ClientRepository, code *CodeInteractor) AuthenticateService {
	return func(ctx context.Context, req *AuthenticateRequest) (*AuthenticateResponse, error) {
		return Authenticate(repo, code, req)
	}
}

// TokenService represents the token flow.
type TokenService func(ctx context.Context, req *TokenRequest) (*TokenResponse, error)

// MakeTokenService returns a new TokenService.
func MakeTokenService(clientRepo ClientRepository, codeRepo CodeRepository, signer Signer) TokenService {
	tokenOpts := TokenOptions{
		AccessTokenDuration:  1 * time.Hour,
		RefreshTokenDuration: 2 * time.Hour,
		Issuer:               "go-openid",
		TokenType:            "bearer",
	}
	return func(ctx context.Context, req *TokenRequest) (*TokenResponse, error) {
		return Token(ctx, tokenOpts, clientRepo, codeRepo, signer, req)
	}
}

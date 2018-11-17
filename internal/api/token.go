package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type (
	TokenRequest struct {
		GrantType   string
		Code        string
		RedirectURI string
	}

	TokenResponse struct {
		AccessToken  string
		TokenType    string
		RefreshToken string
		ExpiresIn    int64
		IDToken      string
	}
)

func Token(
	ctx context.Context,
	clientRepo ClientRepository,
	codeRepo CodeRepository,
	signer Signer,
	req *TokenRequest,
) (*TokenResponse, error) {
	// Pre-Work
	if err := ValidateTokenRequest(req); err != nil {
		return nil, err
	}
	clientID, ok := ctx.Value(ContextKeyClientID).(string)
	if !ok || stringIsEmpty(clientID) {
		return nil, errors.New("client_id is required")
	}
	clientSecret, ok := ctx.Value(ContextKeyClientSecret).(string)
	if !ok || stringIsEmpty(clientSecret) {
		return nil, errors.New("client_secret is required")
	}
	// Do Work.
	client, err := ValidateClientCredentials(clientRepo, clientID, clientSecret)
	if err != nil {
		return nil, err
	}
	if !URIs(client.RedirectURIs).Contains(req.RedirectURI) {
		return nil, errors.New("redirect_uri is invalid")
	}
	// Fetch first , then validate.
	code, err := GetCode(codeRepo, req.Code)
	if err != nil {
		return nil, err
	}
	if code.HasExpired() {
		return nil, errors.New("code is invalid")
	}
	// Get this from the session through the context.
	sub, ok := ctx.Value(ContextKeySubject).(string)
	if !ok || stringIsEmpty(sub) {
		return nil, errors.New("subject is required")
	}

	// The timestamp is passed in through the controller.
	now, ok := ctx.Value(ContextKeyTimestamp).(time.Time)
	if !ok {
		now = time.Now().UTC()
	}

	defaultClaims := jwt.StandardClaims{
		IssuedAt: now.Unix(),
		Subject:  sub,
	}

	accessTokenClaims := defaultClaims
	accessTokenClaims.ExpiresAt = now.Add(2 * time.Hour).Unix()

	refreshTokenClaims := defaultClaims
	refreshTokenClaims.ExpiresAt = now.Add(24 * time.Hour).Unix()

	accessToken, err := signer.Sign(accessTokenClaims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := signer.Sign(refreshTokenClaims)
	if err != nil {
		return nil, err
	}
	claims := NewIDToken()
	idToken, err := signer.Sign(claims)
	if err != nil {
		return nil, err
	}
	return &TokenResponse{
		ExpiresIn:    int64((2 * time.Hour).Seconds()),
		TokenType:    "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IDToken:      idToken,
	}, nil
}

func ValidateTokenRequest(req *TokenRequest) error {
	// Validate required fields.
	if stringIsEmpty(req.Code) {
		return errors.New("code is required")
	}
	if stringIsEmpty(req.GrantType) {
		return errors.New("grant_type is required")
	}
	if stringIsEmpty(req.RedirectURI) {
		return errors.New("redirect_uri is required")
	}
	// Validate type.
	if req.GrantType != "authorization_code" {
		return fmt.Errorf(`grant_type "%s" is invalid`, req.GrantType)
	}
	// Another option is to create the URI type with a validate method.
	if err := ValidateURI(req.RedirectURI); err != nil {
		return fmt.Errorf(`"%s" is not a valid redirect_uri`, req.RedirectURI)
	}
	return nil
}

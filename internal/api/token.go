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

	TokenResponseBuilder interface {
		Build(*TokenRequest) (*TokenResponse, error)
	}
)

func NewTokenResponse(accessToken, refreshToken, idToken string, durationInSeconds int64) *TokenResponse {
	return &TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		RefreshToken: refreshToken,
		ExpiresIn:    durationInSeconds,
		IDToken:      idToken,
	}
}

func Token(
	ctx context.Context,
	clientRepo ClientRepository,
	codeRepo CodeRepository,
	req *TokenRequest,
	responseBuilder TokenResponseBuilder,
) (*TokenResponse, error) {
	// Pre-Work
	if err := ValidateTokenRequest(req); err != nil {
		return nil, err
	}

	clientID, ok := ctx.Value(ContextKeyClientID).(string)
	if !ok || stringIsEmpty(clientID) {
		return nil, errors.New("client_id is required")
	}
	clientSecret, ok := ctx.Value(ContextKeyClientID).(string)
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
	// if code.HasExpired(time.Now())
	if code.HasExpired() {
		return nil, errors.New("code is invalid")
	}
	// if err := ValidateCodeExpiration(code); err != nil {
	//         return err
	// }
	// Post-Work
	// b := builder.TokenResponse()
	// b.Build()
	res, err := responseBuilder.Build(req)
	return res, err
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
	}
	if err := ValidateURI(req.RedirectURI); err != nil {
		return fmt.Errorf(`"%s" is not a valid redirect_uri`, req.RedirectURI)
	}
	return nil
}

// This is actually template design pattern, not builder.
type tokenResponseBuilder struct {
	// The fields here can be set to private to prevent others from modifying it.
	// options TokenResponseBuilderOptions

	// How long before the token expire.
	durationInSeconds int64
	signingKey        []byte
}

func NewTokenResponseBuilder(durationInSeconds int64, signingKey []byte) *tokenResponseBuilder {
	return &tokenResponseBuilder{
		durationInSeconds: durationInSeconds,
		signingKey:        signingKey,
	}
}

// Should not accept a request from build.
func (t *tokenResponseBuilder) Build(req *TokenRequest) (*TokenResponse, error) {
	var accessToken, refreshToken, idToken string
	var err error
	sub := "subject"
	accessToken, err = t.provideAccessToken(sub)
	if err != nil {
		return nil, err
	}
	refreshToken, err = t.provideRefreshToken(sub)
	if err != nil {
		return nil, err
	}
	idToken, err = t.provideIDToken()
	if err != nil {
		return nil, err
	}
	return NewTokenResponse(accessToken, refreshToken, idToken, t.durationInSeconds), nil
}

func (t *tokenResponseBuilder) provideAccessToken(sub string) (string, error) {
	now := time.Now().UTC()
	exp := now.Add(time.Duration(t.durationInSeconds) * time.Second).Unix()
	iat := now.Unix()
	claims := &jwt.StandardClaims{
		Audience:  "https://server.example.com",
		ExpiresAt: exp,
		IssuedAt:  iat,
		Issuer:    "openid",
		Subject:   sub, // UserID
	}
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// return token.SignedString(t.signingKey)
	return signJWT(t.signingKey, claims)
}

func (t *tokenResponseBuilder) provideRefreshToken(sub string) (string, error) {
	now := time.Now().UTC()
	exp := now.Add(time.Duration(2*t.durationInSeconds) * time.Second).Unix()
	iat := now.Unix()
	claims := &jwt.StandardClaims{
		Audience:  "https://server.example.com",
		ExpiresAt: exp,
		IssuedAt:  iat,
		Issuer:    "openid",
		Subject:   sub, // UserID
	}
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// return token.SignedString(t.signingKey)
	return signJWT(t.signingKey, claims)
}

func (t *tokenResponseBuilder) provideIDToken() (string, error) {
	claims := NewIDToken()
	return signJWT(t.signingKey, claims)
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// return token.SignedString(t.signingKey)
}

func signJWT(key []byte, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

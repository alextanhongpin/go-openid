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

	TokenResponseBuilder struct {
		defaults TokenResponse
		override func(t *TokenResponse)
	}

	TokenResponseFactory interface {
		Build(...TokenResponseModifier) (*TokenResponse, error)
		SetOverride(TokenResponseModifier)
	}

	TokenResponseModifier func(t *TokenResponse) error
)

func Token(
	ctx context.Context,
	clientRepo ClientRepository,
	codeRepo CodeRepository,
	responseFactory TokenResponseFactory,
	claimFactory ClaimFactory,
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
	if code.HasExpired() {
		return nil, errors.New("code is invalid")
	}
	// Get this from the session through the context.
	sub := "hello"

	// The timestamp is passed in through the controller.
	now, ok := ctx.Value(ContextKeyTimestamp).(time.Time)
	if !ok {
		now = time.Now().UTC()
	}

	issuedAtModifier := makeIssuedAtModifier(now)
	subjectModifier := makeSubjectModifier(sub)
	expiresIn := 2 * time.Hour
	accessToken := claimFactory.Build(
		issuedAtModifier,
		subjectModifier,
		makeExpireAtModifier(now, expiresIn),
	)
	refreshToken := claimFactory.Build(
		issuedAtModifier,
		subjectModifier,
		makeExpireAtModifier(now, 24*time.Hour),
	)
	res, err := responseFactory.Build(
		makeAccessTokenExpiresIn(int64(expiresIn.Seconds())),
		makeAccessTokenModifier(signer, accessToken),
		makeRefreshTokenModifier(signer, refreshToken),
		makeIDTokenModifier(signer),
	)
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
		return fmt.Errorf(`grant_type "%s" is invalid`, req.GrantType)
	}
	if err := ValidateURI(req.RedirectURI); err != nil {
		return fmt.Errorf(`"%s" is not a valid redirect_uri`, req.RedirectURI)
	}
	return nil
}

func NewTokenResponseBuilder(defaults TokenResponse) *TokenResponseBuilder {
	return &TokenResponseBuilder{
		defaults: defaults,
		override: func(t *TokenResponse) {},
	}
}
func (t *TokenResponseBuilder) SetOverride(override func(t *TokenResponse)) {
	t.override = override
}
func (t *TokenResponseBuilder) SetExpiresIn(durationInSeconds int64) {
	t.defaults.ExpiresIn = durationInSeconds
}
func (t *TokenResponseBuilder) SetAccessToken(token string) {
	t.defaults.AccessToken = token
}
func (t *TokenResponseBuilder) SetRefreshToken(token string) {
	t.defaults.RefreshToken = token
}
func (t *TokenResponseBuilder) SetIDToken(token string) {
	t.defaults.IDToken = token
}
func (t *TokenResponseBuilder) Build() *TokenResponse {
	result := t.defaults
	if t.override != nil {
		t.override(&result)
	}
	return &result
}

// This is actually template design pattern, not builder.
type tokenResponseFactory struct {
	// The fields here can be set to private to prevent others from modifying it.
	// options TokenResponseFactoryOptions

	// How long before the token expire.
	defaults TokenResponse
	override TokenResponseModifier
}

func NewTokenResponseFactory() *tokenResponseFactory {
	return &tokenResponseFactory{
		defaults: TokenResponse{
			TokenType: "Bearer",
			ExpiresIn: 3600,
		},
	}
}
func (t *tokenResponseFactory) SetOverride(override TokenResponseModifier) {
	t.override = override
}

// Should not accept a request from build.
func (t *tokenResponseFactory) Build(modifiers ...TokenResponseModifier) (*TokenResponse, error) {
	var err error
	result := t.defaults
	for _, modifier := range modifiers {
		err = modifier(&result)
		if err != nil {
			return nil, err
		}
	}
	if t.override != nil {
		err = t.override(&result)
	}
	return &result, err
}

func makeRefreshTokenModifier(signer Signer, claims jwt.Claims) TokenResponseModifier {
	return func(t *TokenResponse) error {
		var err error
		t.RefreshToken, err = signer.Sign(claims)
		return err
	}
}
func makeAccessTokenModifier(signer Signer, claims jwt.Claims) TokenResponseModifier {
	return func(t *TokenResponse) error {
		var err error
		t.AccessToken, err = signer.Sign(claims)
		return err
	}
}
func makeIDTokenModifier(signer Signer) TokenResponseModifier {
	return func(t *TokenResponse) error {
		var err error
		claims := NewIDToken()
		t.IDToken, err = signer.Sign(claims)
		return err
	}
}

func makeAccessTokenExpiresIn(timestamp int64) TokenResponseModifier {
	return func(t *TokenResponse) error {
		t.ExpiresIn = timestamp
		return nil
	}
}

package core

import (
	"context"
	"errors"
	"time"

	"github.com/alextanhongpin/go-openid"
	database "github.com/alextanhongpin/go-openid/internal/database"
	"github.com/alextanhongpin/go-openid/pkg/authheader"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
	"github.com/alextanhongpin/go-openid/repository"
	"github.com/asaskevich/govalidator"
	jwt "github.com/dgrijalva/jwt-go"
)

type modelImpl struct {
	code   repository.Code
	client repository.Client
	user   repository.User
}

func NewModel() *modelImpl {
	return &modelImpl{
		code:   database.NewCodeKV(),
		client: database.NewClientKV(),
		user:   database.NewUserKV(),
	}
}

// ValidateAuthnRequest validates the required fields for the authentication
// request.
func (m *modelImpl) ValidateAuthnRequest(req *oidc.AuthenticationRequest) error {
	var (
		clientID     = req.ClientID
		redirectURI  = req.RedirectURI
		prompt       = req.GetPrompt()
		responseType = req.GetResponseType()
		scope        = req.GetScope()
	)
	// Scope cannot be none, and it should have at least an openid scope.
	if scope.Is(oidc.ScopeNone) || !scope.Has(oidc.ScopeOpenID) {
		return errors.New("scope required")
	}

	// ResponseType cannot be none, and should have "code" for
	// authorization code flow.
	if responseType.Is(oidc.ResponseTypeNone) || !responseType.Has(oidc.ResponseTypeCode) {
		return errors.New("response_type required")
	}

	if clientID == "" {
		return errors.New("client_id required")
	}

	if redirectURI == "" {
		return errors.New("redirect_uri required")
	}

	if !govalidator.IsURL(redirectURI) {
		return errors.New("redirect_uri invalid")
	}

	// If prompt is "none", it cannot have other values.
	if prompt.Has(oidc.PromptNone) && prompt.Has(oidc.PromptLogin|oidc.PromptConsent|oidc.PromptSelectAccount) {
		return errors.New("prompt none may not contain other values")
	}
	return nil
}

// ValidateAuthnUser validates the authentication request with the user data in
// the database.
func (m *modelImpl) ValidateAuthnUser(ctx context.Context, req *oidc.AuthenticationRequest) error {
	userID, ok := oidc.GetUserIDContextKey(ctx)
	if !ok {
		return errors.New("user_id missing")
	}
	user, err := m.user.Get(userID)
	if err != nil {
		return err
	}
	if time.Since(time.Unix(user.Profile.UpdatedAt, 0)) > time.Duration(req.MaxAge) {
		// TODO: Must re-authenticate.
		return errors.New("re-authentication required")
	}

	if req.LoginHint == user.Email.Email {
		// user.Email.Email
	}
	return nil
}

// ValidateAuthnClient validates the provided client request with the client
// data in the storage.
func (m *modelImpl) ValidateAuthnClient(req *oidc.AuthenticationRequest) error {
	var (
		clientID    = req.ClientID
		redirectURI = req.RedirectURI
	)
	client, err := m.client.Get(clientID)
	if err != nil {
		return err
	}
	if !client.GetRedirectURIs().Contains(redirectURI) {
		return errors.New("redirect_uri incorrect")
	}
	return nil
}

// NewCode returns a new code.
func (m *modelImpl) NewCode() string {
	c := crypto.NewXID()
	code := oidc.NewCode(c)
	m.code.Put(c, code)
	return c
}

func (m *modelImpl) ValidateClientAuthHeader(authorization string) (*oidc.Client, error) {
	token, err := authheader.Basic(authorization)
	if err != nil {
		return nil, err
	}
	clientID, clientSecret, err := authheader.DecodeBase64(token)
	if err != nil {
		return nil, err
	}
	return m.client.GetByCredentials(clientID, clientSecret)
}

func (m *modelImpl) ProvideToken(userID string, duration time.Duration) (string, error) {
	var (
		aud = "https://server.example.com/token"
		sub = userID
		iss = userID
		iat = time.Now().UTC()
		exp = iat.Add(duration)

		key = []byte("access_token_secret")
	)
	claims := crypto.NewStandardClaims(aud, sub, iss, iat.Unix(), exp.Unix())
	return crypto.NewJWT(key, claims)
}

func (m *modelImpl) ProvideIDToken(userID string) (string, error) {
	user, err := m.user.Get(userID)
	if err != nil {
		return "", err
	}
	idToken := user.ToIDToken()
	var (
		now = time.Now().UTC()
		aud = "https://server.example.com/token"
		sub = userID
		iss = userID // Should be client id?
		iat = now
		id  = crypto.NewXID()
		nbf = now
		exp = now.Add(2 * time.Hour)
	)
	idToken.StandardClaims = jwt.StandardClaims{
		Audience:  aud,
		ExpiresAt: exp.Unix(),
		Id:        id,
		IssuedAt:  iat.Unix(),
		Issuer:    iss,
		NotBefore: nbf.Unix(),
		Subject:   sub,
	}
	return idToken.SignHS256([]byte("id_token_key"))
}

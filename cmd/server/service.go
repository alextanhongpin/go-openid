package main

import (
	"context"
	"errors"
	"strings"
	"time"

	oidc "github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/internal/database"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
)

// Service represents the interface for the services available for OpenID Connect Protocol.
type Service interface {
	Authorize(context.Context, *oidc.AuthorizationRequest) (*oidc.AuthorizationResponse, error)
	Token(context.Context, *oidc.AccessTokenRequest) (*oidc.AccessTokenResponse, error)
	RegisterClient(context.Context, *oidc.ClientRegistrationRequest) (*oidc.ClientRegistrationResponse, error)
	Client(context.Context, string) (*oidc.Client, error)
	UserInfo(context.Context, string) (*oidc.StandardClaims, error)
	ValidateClient(clientID, clientSecret string) error
	ParseJWT(token string) (*oidc.Claims, error)
	// RegisterUser
	// Authenticate
}

// ServiceImpl fulfils the OIDService interface.
type ServiceImpl struct {
	crypto crypto.Crypto
	db     *database.Database
}

// NewService returns a pointer to a new service.
func NewService(db *database.Database, c crypto.Crypto) *ServiceImpl {
	return &ServiceImpl{
		crypto: c,
		db:     db,
	}
}

func (s *ServiceImpl) Authorize(ctx context.Context, req *oidc.AuthorizationRequest) (*oidc.AuthorizationResponse, error) {
	if err := req.Validate(); err != nil {
		if req.State != "" {
			if e, ok := err.(*oidc.ErrorJSON); ok {
				e.SetState(req.State)
				return nil, e
			}
		}
		return nil, err
	}

	var (
		cid   = req.ClientID
		ruri  = req.RedirectURI
		state = req.State
	)

	_, err := s.validateClient(cid, ruri)
	if err != nil {
		if e, ok := err.(*oidc.ErrorJSON); ok {
			e.SetState(state)
			return nil, e
		}
		return nil, err
	}

	code := s.newCode(cid)

	// Return response
	return &oidc.AuthorizationResponse{
		State: req.State,
		Code:  code.Code,
	}, nil
}

func (s *ServiceImpl) Token(ctx context.Context, req *oidc.AccessTokenRequest) (*oidc.AccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	var (
		cid  = req.ClientID
		ruri = req.RedirectURI
		code = req.Code
	)

	client, err := s.validateClient(cid, ruri)
	if err != nil {
		return nil, err
	}

	if err := s.validateCode(cid, code); err != nil {
		return nil, err
	}

	var (
		aud = client.ClientName
		sub = client.ClientID
		iss = defaultIssuer
		dur = defaultDuration
	)

	accessToken, err := s.crypto.NewJWT(aud, sub, iss, dur)
	if err != nil {
		return nil, err
	}

	// Make refresh token last longer than the average token.
	refreshToken, err := s.crypto.NewJWT(aud, sub, iss, dur*4)
	if err != nil {
		return nil, err
	}

	// Finalize the response and return the access token.
	return &oidc.AccessTokenResponse{
		AccessToken:  accessToken,
		TokenType:    "bearer",
		ExpiresIn:    int64(defaultDuration.Seconds()),
		RefreshToken: refreshToken,
		IDToken:      "",
	}, nil
}

// RegisterClient represents the dynamic client registration at the connect/register endpoint
func (s *ServiceImpl) RegisterClient(ctx context.Context, req *oidc.ClientRegistrationRequest) (*oidc.ClientRegistrationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Check if client is already registered
	// If client is already registered, return err
	if _, exist := s.db.Client.Get(req.ClientName); exist {
		return nil, errors.New("client is already registered")
	}

	// If the client is not registered, create a new client
	client, err := s.newClient(req)
	if err != nil {
		return nil, err
	}

	// Save the client to the storage
	s.db.Client.Put(req.ClientName, client)

	return client.ClientPrivate, nil
}

// Client returns the client from the given client id.
func (s *ServiceImpl) Client(ctx context.Context, id string) (*oidc.Client, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, oidc.UnauthorizedClient.JSON()
	}
	if c := s.db.Client.GetByID(id); c != nil {
		return c, nil
	}
	return nil, oidc.UnauthorizedClient.JSON()
}

// UserInfo returns the user info from the given id.
func (s *ServiceImpl) UserInfo(ctx context.Context, id string) (*oidc.StandardClaims, error) {
	user, exist := s.db.User.Get(id)
	if !exist || user == nil {
		return nil, oidc.AccessDenied.JSON()
	}

	return user, nil
}

// RefreshToken represents the refresh token flow.
func (s *ServiceImpl) RefreshToken(ctx context.Context, req *oidc.RefreshTokenRequest) (*oidc.RefreshTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	client := s.db.Client.GetByID(req.ClientID)
	if client == nil {
		return nil, oidc.UnauthorizedClient.JSON()
	}

	if client.ClientPrivate.ClientSecret != req.ClientSecret {
		return nil, oidc.UnauthorizedClient.JSON()
	}

	var (
		aud = client.ClientName
		dur = defaultDuration
		iss = defaultIssuer
		sub = client.ClientID
	)

	accessToken, err := s.crypto.NewJWT(aud, sub, iss, dur)
	if err != nil {
		return nil, err
	}

	return &oidc.RefreshTokenResponse{
		AccessToken: accessToken,
		// TODO: Return the expire time and also ID token.
	}, nil
}

// ValidateClient attempts to validate the provided client id and client secret
// by checking the database.
func (s *ServiceImpl) ValidateClient(clientID, clientSecret string) error {
	client := s.db.Client.GetByIDAndSecret(clientID, clientSecret)
	if client == nil {
		return oidc.UnauthorizedClient.JSON()
	}
	return nil
}

// ParseJWT takes a jwt token and return the decoded claims.
func (s *ServiceImpl) ParseJWT(token string) (*oidc.Claims, error) {
	return s.crypto.ParseJWT(token)
}

func (s *ServiceImpl) newClient(req *oidc.ClientPublic) (*oidc.Client, error) {
	var (
		aud = req.ClientName
		dur = time.Hour * 24 * 30 // 1 month.
		iss = defaultIssuer
		sub = req.ClientName

		now = time.Now().UTC()
	)

	token, err := s.crypto.NewJWT(aud, iss, sub, dur)
	if err != nil {
		return nil, err
	}

	clientPrivate := &oidc.ClientPrivate{
		ClientID:                s.crypto.Code(),
		ClientSecret:            s.crypto.UUID(),
		RegistrationAccessToken: token,
		RegistrationClientURI:   "",
		ClientIDIssuedAt:        now.Unix(),
		ClientSecretExpiresAt:   now.Add(dur).Unix(),
	}

	return &oidc.Client{
		ClientPublic:  req,
		ClientPrivate: clientPrivate,
	}, nil
}

func (s *ServiceImpl) validateClient(cid, ruri string) (*oidc.Client, error) {
	if len(cid) == 0 || len(ruri) == 0 {
		return nil, errors.New("forbidden request")
	}

	client := s.db.Client.GetByID(cid)
	if client == nil {
		return nil, oidc.InvalidClientMetadata.JSON()
	}
	if match := client.RedirectURIs.Contains(ruri); !match {
		return nil, oidc.InvalidRedirectURI.JSON()
	}

	return client, nil
}

func (s *ServiceImpl) newCode(cid string) *oidc.Code {
	// Delete existing code.
	if _, exist := s.db.Code.Get(cid); exist {
		s.db.Code.Delete(cid)
	}

	// Create new code and store it.
	code := oidc.NewCode(s.crypto.Code())
	s.db.Code.Put(cid, code)
	return code
}

func (s *ServiceImpl) validateCode(cid, code string) error {
	// Check if the code exists, and it matches the code provided
	c, exist := s.db.Code.Get(cid)
	if !exist || c == nil || c.Code != code {
		return oidc.AccessDenied.JSON()
	}

	// If the code is valid, but expired, delete them
	if c.Expired() {
		s.db.Code.Delete(cid)
		return oidc.AccessDenied.JSON()
	}

	// If code matches, then delete the code from the storage
	s.db.Code.Delete(cid)
	return nil
}

package main

import (
	"context"
	"errors"
	"time"

	oidc "github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
)

var (
	defaultIssuer        = "go-openid"
	defaultDuration      = time.Hour
	defaultJWTSigningKey = "secret"
)

// Service represents the interface for the services available for OpenID Connect Protocol.
type Service interface {
	Authorize(context.Context, *oidc.AuthorizationRequest) (*oidc.AuthorizationResponse, *oidc.AuthorizationError)
	Token(context.Context, *oidc.AccessTokenRequest) (*oidc.AccessTokenResponse, error)
	RegisterClient(context.Context, *oidc.ClientRegistrationRequest) (*oidc.ClientRegistrationResponse, error)
	UserInfo(context.Context, string) (*User, error)
	// RegisterUser
	// Authenticate
}

// ServiceImpl fulfils the OIDService interface.
type ServiceImpl struct {
	crypto crypto.Crypto
	db     *Database
}

// NewService returns a pointer to a new service.
func NewService(db *Database, c crypto.Crypto) *ServiceImpl {
	if db == nil {
		db = NewDatabase()
	}
	if c == nil {
		c = crypto.New(defaultJWTSigningKey)
	}
	return &ServiceImpl{
		crypto: c,
		db:     db,
	}
}

func (s *ServiceImpl) newClient(req *oidc.ClientPublic) (*oidc.Client, error) {
	var (
		dur = time.Hour * 24 * 30
		aud = req.ClientName
		iss = defaultIssuer
		sub = req.ClientName
	)

	token, err := s.crypto.NewJWT(aud, iss, sub, dur)
	if err != nil {
		return nil, err
	}
	now := time.Now()
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
		return nil, errors.New("client is not authorized")
	}
	if match := client.RedirectURIs.Contains(ruri); !match {
		return nil, errors.New("one or more redirect uris are incorrect")
	}
	return client, nil
}

func (s *ServiceImpl) newCode(cid string) *oidc.Code {
	if _, exist := s.db.Code.Get(cid); exist {
		s.db.Code.Delete(cid)
	}
	code := oidc.NewCode(s.crypto.Code())
	s.db.Code.Put(cid, code)
	return code
}

func (s *ServiceImpl) validateCode(cid, code string) error {
	// Check if the code exists, and it matches the code provided
	c, exist := s.db.Code.Get(cid)
	if !exist || c == nil || c.Code != code {
		return oidc.ErrForbidden
	}

	// If the code is valid, but expired, delete them
	if c.Expired() {
		s.db.Code.Delete(cid)
		return errors.New("expired code")
	}

	// If code matches, then delete the code from the storage
	s.db.Code.Delete(cid)
	return nil
}

func (s *ServiceImpl) Authorize(ctx context.Context, req *oidc.AuthorizationRequest) (*oidc.AuthorizationResponse, *oidc.AuthorizationError) {
	if err := req.Validate(); err != nil {
		return nil, &oidc.AuthorizationError{
			Error:            oidc.ErrForbidden.Error(),
			ErrorDescription: "",
			ErrorURI:         "",
			State:            req.State,
		}
	}
	var (
		cid   = req.ClientID
		ruri  = req.RedirectURI
		state = req.State
	)

	_, err := s.validateClient(cid, ruri)
	if err != nil {
		return nil, &oidc.AuthorizationError{
			Error:            err.Error(),
			ErrorDescription: "",
			ErrorURI:         "",
			State:            state,
		}
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

	refreshToken, err := s.crypto.NewJWT(aud, sub, iss, dur*4)
	if err != nil {
		return nil, err
	}

	// Finalize the response and return the access token
	return &oidc.AccessTokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(defaultDuration.Seconds()),
		RefreshToken: refreshToken,
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

func (s *ServiceImpl) UserInfo(ctx context.Context, id string) (*User, error) {
	user, exist := s.db.User.Get(id)
	if !exist || user == nil {
		return nil, errors.New("forbidden access")
	}
	return user, nil
}

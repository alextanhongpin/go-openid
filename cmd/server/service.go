package main

import (
	"context"
	"errors"

	oidc "github.com/alextanhongpin/go-openid"
	"github.com/alextanhongpin/go-openid/pkg/crypto"
)

// Service represents the interface for the services available for OpenID Connect Protocol.
type Service interface {
	Authorize(context.Context, *oidc.AuthorizationRequest) (*oidc.AuthorizationResponse, *oidc.AuthorizationError)
	Token(context.Context, *oidc.AccessTokenRequest) (*oidc.AccessTokenResponse, error)
	RegisterClient(context.Context, *oidc.ClientRegistrationRequest) (*oidc.ClientRegistrationResponse, error)
	UserInfo(context.Context, string) (*User, error)
	// RegisterUser
	// Authenticate
	GenerateCode() string
	NewJWT(aud, iss, sub string) (string, error)
}

// tokenGenerator represents the function to generate token.
type tokenGenerator func() string

type Client interface {
	New(*oidc.ClientRegistrationRequest) *oidc.Client
}

// ServiceImpl fulfils the OIDService interface.
type ServiceImpl struct {
	crypto               crypto.Crypto
	db                   *Database
	generateCode         tokenGenerator
	generateAccessToken  tokenGenerator
	generateRefreshToken tokenGenerator
	newClient            func(*oidc.ClientRegistrationRequest) *oidc.Client
}

// NewService returns a pointer to a new service.
func NewService(db *Database, gc tokenGenerator, gat tokenGenerator, grt tokenGenerator) *ServiceImpl {
	if db == nil {
		db = NewDatabase()
	}
	return &ServiceImpl{
		db:                   db,
		generateCode:         gc,
		generateAccessToken:  gat,
		generateRefreshToken: grt,
		newClient: func(req *oidc.ClientRegistrationRequest) *oidc.Client {
			return oidc.NewClient(req)
		},
	}
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
	cid := req.ClientID

	// Check if client exist
	client := s.db.Client.GetByID(cid)
	if client == nil {
		return nil, &oidc.AuthorizationError{
			Error:            oidc.ErrForbidden.Error(),
			ErrorDescription: "",
			ErrorURI:         "",
			State:            req.State,
		}
	}

	// Check if redirect uri is correct
	if match := client.RedirectURIs.Contains(req.RedirectURI); !match {
		// TODO: Return the error query string in the redirect uri
		return nil, &oidc.AuthorizationError{
			Error:            oidc.ErrForbidden.Error(),
			ErrorDescription: "one or more redirect uris are incorrect",
			ErrorURI:         "",
			State:            req.State,
		}

	}

	// Check if client has code, if yes, remove existing ones and return a new one
	if _, exist := s.db.Code.Get(cid); exist {
		s.db.Code.Delete(cid)
	}
	newCode := oidc.NewCode(s.generateCode())

	// Set the code in the storage
	s.db.Code.Put(cid, newCode)

	// Return response
	return &oidc.AuthorizationResponse{
		State: req.State,
		Code:  newCode.Code,
	}, nil
}

func (s *ServiceImpl) Token(ctx context.Context, req *oidc.AccessTokenRequest) (*oidc.AccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	cid := req.ClientID
	// Check if the client exists
	client := s.db.Client.GetByID(cid)
	if client == nil {
		return nil, oidc.ErrForbidden
	}
	// Check if redirect uri is correct
	if match := client.RedirectURIs.Contains(req.RedirectURI); !match {
		return nil, errors.New("one or more redirect uri is incorrect")
	}

	// Check if the code exists, and it matches the code provided
	code, exist := s.db.Code.Get(cid)
	if !exist || code.Code != req.Code {
		return nil, oidc.ErrForbidden
	}

	// If the code is valid, but expired, delete them
	if code.Expired() {
		s.db.Code.Delete(cid)
		return nil, errors.New("expired code")
	}

	// If code matches, then delete the code from the storage
	s.db.Code.Delete(cid)

	// accessToken := s.NewJWT(client.ClientName, s.Config.issuer, client.ClientID)

	// Finalize the response and return the access token
	return &oidc.AccessTokenResponse{
		AccessToken:  s.generateAccessToken(),
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: s.generateRefreshToken(),
	}, nil
}

// RegisterClient represents the dynamic client registration at the connect/register endpoint
func (s *ServiceImpl) RegisterClient(ctx context.Context, req *oidc.ClientRegistrationRequest) (*oidc.ClientRegistrationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	// Check if client is already registered
	if _, exist := s.db.Client.Get(req.ClientName); exist {
		return nil, errors.New("client is already registered")
	}

	// If client is already registered, return err
	// If the client is not registered, create a new client
	client := s.newClient(req)

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

package main

import (
	"context"
	"errors"
	"strings"

	oidc "github.com/alextanhongpin/go-openid"
)

// OIDService represents the interface for the services available for OpenID Protocol
type OIDService interface {
	Authorize(context.Context, *oidc.AuthorizationRequest) (*oidc.AuthorizationResponse, error)
	Token(context.Context, *oidc.AccessTokenRequest) (*oidc.AccessTokenResponse, error)
	// RegisterClient
	// RegisterUser
	// Authenticate
}

type tokenGenerator func() string

// Service fulfils the OIDService interface.
type Service struct {
	db      *Database
	genCode tokenGenerator
	genAT   tokenGenerator
	genRT   tokenGenerator
}

// NewService returns a pointer to a new service.
func NewService(db *Database, codeGen tokenGenerator, atGen tokenGenerator, rtGen tokenGenerator) *Service {
	if db == nil {
		db = NewDatabase()
	}
	return &Service{
		db:      db,
		genCode: codeGen,
		genAT:   atGen,
		genRT:   rtGen,
	}
}

func validateRedirectURIs(uris []string, uri string) bool {
	for _, u := range uris {
		if strings.Compare(u, uri) == 0 {
			return true
		}
	}
	return false
}

func (s *Service) Authorize(ctx context.Context, req *oidc.AuthorizationRequest) (*oidc.AuthorizationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	cid := req.ClientID

	// Check if client exist
	client, exist := s.db.Client.Get(cid)
	if !exist || client == nil {
		return nil, oidc.ErrForbidden
	}

	// Check if redirect uri is correct
	if match := validateRedirectURIs(client.RedirectURIs, req.RedirectURI); !match {
		return nil, errors.New("one or more redirect uri is incorrect")
	}

	// Check if client has code, if yes, remove existing ones and return a new one
	if _, exist = s.db.Code.Get(cid); exist {
		s.db.Code.Delete(cid)
	}
	newCode := oidc.NewCode(s.genCode())

	// Set the code in the storage
	s.db.Code.Put(cid, newCode)

	// Return response
	return &oidc.AuthorizationResponse{
		State: req.State,
		Code:  newCode.Code,
	}, nil
}

func (s *Service) Token(ctx context.Context, req *oidc.AccessTokenRequest) (*oidc.AccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	cid := req.ClientID
	// Check if the client exists
	client, exist := s.db.Client.Get(cid)
	if !exist || client == nil {
		return nil, oidc.ErrForbidden
	}

	// Check if redirect uri is correct
	if match := validateRedirectURIs(client.RedirectURIs, req.RedirectURI); !match {
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

	// Finalize the response and return the access token
	return &oidc.AccessTokenResponse{
		AccessToken:  s.genAT(),
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: s.genRT(),
	}, nil
}

func (s *Service) RegisterClient(ctx context.Context, req *oidc.ClientRegistrationRequest) (*oidc.ClientRegistrationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	// Check if client is already registered
	// s.db.Client.Get(req.?)

	// If client is already registered, return err
	// If the client is not registered, create a new client
	// client := NewClient(clientID, clientSecret)

	// Save the client to the storage
	// s.db.Client.Put(clientID, client)

	return nil, nil
}

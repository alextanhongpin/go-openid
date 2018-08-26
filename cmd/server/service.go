package main

import (
	"context"
	"errors"
	"strings"
	"time"

	openid "github.com/alextanhongpin/go-openid"
)

// OIDService represents the interface for the services available for OpenID Protocol
type OIDService interface {
	Authorize(context.Context, *openid.AuthorizationRequest) (*openid.AuthorizationResponse, error)
	Token(context.Context, *openid.AccessTokenRequest) (*openid.AccessTokenResponse, error)
	// GenerateCode() (string, error)
	// GenerateAccessToken(string) (string, error)
}

type tokenGenerator func() string

// Service fulfils the OIDService interface
type Service struct {
	db      *Database
	genCode tokenGenerator
	genAT   tokenGenerator
	genRT   tokenGenerator
}

func (s *Service) Authorize(ctx context.Context, req *openid.AuthorizationRequest) (*openid.AuthorizationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	cid := req.ClientID

	// Check if client exist
	client, exist := s.db.Client.Get(cid)
	if !exist || client == nil {
		return nil, openid.ErrForbidden
	}

	// Check if redirect uri is correct
	if !strings.Contains(client.RedirectURIs, req.RedirectURI) {
		return nil, errors.New("one or more redirect uri is incorrect")
	}

	// Check if client has code, if yes, remove existing ones and return a new one
	_, exist = s.db.Code.Get(cid)
	if exist {
		s.db.Code.Delete(cid)
	}
	code := s.genCode()
	newCode := &openid.Code{
		Code:      code,
		CreatedAt: time.Now(),
		TTL:       10 * time.Minute,
	}
	// Set the code in the storage
	s.db.Code.Put(cid, newCode)

	// Return response
	return &openid.AuthorizationResponse{
		State: req.State,
		Code:  code,
	}, nil
}

func (s *Service) Token(ctx context.Context, req *openid.AccessTokenRequest) (*openid.AccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	cid := req.ClientID
	// Check if the client exists
	client, exist := s.db.Client.Get(cid)
	if !exist || client == nil {
		return nil, openid.ErrForbidden
	}

	// Check if redirect uri is correct
	if !strings.Contains(client.RedirectURIs, req.RedirectURI) {
		return nil, errors.New("one or more redirect uri is incorrect")
	}
	// Check if the code exists, and it matches the code provided
	code, exist := s.db.Code.Get(cid)
	if !exist || code.Code != req.Code {
		return nil, openid.ErrForbidden
	}
	// If the code is valid, but expired, delete them
	if code.Expired() {
		s.db.Code.Delete(cid)
		return nil, errors.New("expired code")
	}

	// If code matches, then delete the code from the storage
	s.db.Code.Delete(cid)

	// Finalize the response and return the access token
	return &openid.AccessTokenResponse{
		AccessToken:  s.genAT(),
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: s.genRT(),
	}, nil
}

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

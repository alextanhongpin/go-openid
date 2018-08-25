package openid

import (
	"errors"
	"log"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

type Config struct {
	AccessTokenDuration time.Duration
	CodeDuration        time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		AccessTokenDuration: 2 * time.Hour,
		CodeDuration:        10 * time.Minute,
	}
}

type OpenID interface {
	AuthorizationCode(AuthorizationRequest) (*AuthorizationResponse, error)
	Token(AccessTokenRequest) (*AccessTokenResponse, error)
}

type Service struct {
	config      Config
	CodeStore   CodeStore
	ClientStore ClientStore
}

func (s *Service) AuthorizationCode(req AuthorizationRequest) (*AuthorizationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// TODO: Generate the code
	code := s.generateCode()

	// Store the code in the store with the current time
	s.CodeStore.Put(req.ClientID, code)

	return &AuthorizationResponse{
		Code:  code.Code,
		State: req.State,
	}, nil
}

func (s *Service) validateCode(clientID, code string) error {
	c := s.CodeStore.Get(clientID)
	if c.Code != code {
		return errors.New("invalid code")
	}
	if c.Expired() {
		// Remove expired code
		s.CodeStore.Delete(clientID)
		return errors.New("code expired")
	}
	// Delete the authorization code so that it can't be reused
	s.CodeStore.Delete(clientID)
	return nil
}

func (s *Service) Token(req AccessTokenRequest) (*AccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Check if the client exist
	client := s.ClientStore.Get(req.ClientID)
	if client == nil {
		return nil, ErrForbidden
	}

	// Check if client redirect uri is valid
	if !strings.Contains(client.RedirectURIs, req.RedirectURI) {
		return nil, ErrForbidden
	}

	// Check if the token is valid
	if err := s.validateCode(req.ClientID, req.Code); err != nil {
		return nil, err
	}

	// Finalize the response
	return &AccessTokenResponse{
		AccessToken:  s.generateAccessToken(req.ClientID),
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.config.AccessTokenDuration.Seconds()),
		RefreshToken: s.generateRefreshToken(req.ClientID),
	}, nil
}

func (s *Service) generateAccessToken(clientID string) string {
	claims := &jwt.StandardClaims{
		ExpiresAt: int64(s.config.AccessTokenDuration.Seconds()),
		Issuer:    "test",
		Subject:   clientID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(SigningKey)
	if err != nil {
		log.Println(err)
	}
	return ss
}

func (s *Service) generateRefreshToken(clientID string) string {
	return "refresh_token"
}

func (s *Service) generateCode() *Code {
	return &Code{
		Code:      uuid.Must(uuid.NewV4()).String(),
		CreatedAt: time.Now(),
		TTL:       s.config.CodeDuration,
	}
}

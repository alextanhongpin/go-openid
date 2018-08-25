package openid

import (
	"log"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type Config struct {
	AccessTokenDuration int64
}

func DefaultConfig() *Config {
	return &Config{
		AccessTokenDuration: 10000,
	}
}

type Service struct {
	sync.RWMutex
	config           Config
	ClientCache      map[string]struct{}
	CodeCache        map[string]time.Time
	RedirectURICache map[string]string
}

func (s *Service) generateCode() string {
	return ""
}

func (s *Service) CheckCode (clientID string) bool {
	s.RLock()
	duration, ok := s.CodeCache.Get(clientID)
	s.RUnlock()
	return ok && time.Since(duration) < s.config.CodeTTL
}

func (s *Service) AuthorizationCode(req AuthorizationCodeRequest) (*AuthorizationCodeResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	// Check if the client has generated the code before

	code := s.generateCode()

	return &AuthorizationCodeResponse{
		Code:  
		State: req.State,
	}, nil
}
func (s *Service) generateAccessToken() string {
	claims := &jwt.StandardClaims{
		ExpiresAt: s.config.AccessTokenDuration,
		Issuer:    "test",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(SigningKey)
	log.Println(err)

	return ss
}
func (s *Service) RequestAccessToken(req AccessTokenRequest) (*AccessTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	s.RLock()
	_, hasClient := s.ClientCache[req.ClientID]
	codeDuration, hasCode := s.CodeCache[req.Code]
	redirectURI, hasRedirectURI := s.RedirectURICache[req.ClientID]
	s.RUnlock()
	if !hasClient {
		return nil, ErrForbidden
	}
	if !hasCode || time.Since(codeDuration) > 10*time.Minute {
		return nil, ErrForbidden
	}
	if !hasRedirectURI || redirectURI != req.RedirectURI {
		return nil, ErrForbidden
	}

	res := AccessTokenResponse{
		AccessToken:  s.generateAccessToken(),
		TokenType:    "bearer",
		ExpiresIn:    s.config.AccessTokenDuration,
		RefreshToken: s.generateRefreshToken(),
	}
	return nil, nil
}

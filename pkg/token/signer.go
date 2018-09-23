package token

import (
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

type Signer interface {
	NewClaims(opts ...Option) *Claims
	NewJWT(claims *Claims) (string, error)
	ParseJWT(token string) (*Claims, error)
}

// signerImpl represents the signer of the claims and holds the signing key.
type signerImpl struct {
	key           []byte
	defaultClaims *Claims
}

// NewSigner returns the signer and the default claims to be provided.
func NewSigner(key []byte, defaultClaims *Claims) *signerImpl {
	return &signerImpl{
		key:           key,
		defaultClaims: defaultClaims,
	}
}

// NewClaims returns a new claims on top of the provided option.
func (s *signerImpl) NewClaims(opts ...Option) *Claims {
	return NewClaims(s.defaultClaims, opts...)
}

// NewJWT returns a new jwt signed string from the given claims.
func (s *signerImpl) NewJWT(claims *Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.key)
}

// ParseJWT attempts to parse the raw token string and return the claims.
func (s *signerImpl) ParseJWT(token string) (*Claims, error) {
	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})

	if t.Valid {
		if claims, ok := t.Claims.(*jwt.StandardClaims); ok && t.Valid {
			return claims, nil
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		switch {
		case ve.Errors&jwt.ValidationErrorMalformed != 0:
			return nil, errors.New("invalid token")
		case ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0:
			return nil, errors.New("token expired")
		default:
			return nil, fmt.Errorf("token malformed: %s", err.Error())
		}
	}
	return nil, fmt.Errorf("token malformed: %s", err.Error())
}

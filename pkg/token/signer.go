package token

import (
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

// Signer represents the signer of the claims and holds the signing key.
type Signer struct {
	key           []byte
	defaultClaims *Claims
}

// NewSigner returns the signer and the default claims to be provided.
func NewSigner(key []byte, defaultClaims *Claims) *Signer {
	return &Signer{
		key:           key,
		defaultClaims: defaultClaims,
	}
}

// NewClaims returns a new claims on top of the provided option.
func (s *Signer) NewClaims(opts ...Option) *Claims {
	return NewClaims(s.defaultClaims, opts...)
}

// NewJWT returns a new jwt signed string from the given claims.
func (s *Signer) NewJWT(claims *Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.key)
}

// ParseJWT attempts to parse the raw token string and return the claims.
func (s *Signer) ParseJWT(token string) (*Claims, error) {
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

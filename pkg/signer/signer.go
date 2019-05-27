package signer

import (
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

type Signer struct {
	secret []byte
}

func NewSigner(secret string) *Signer {
	return &Signer{[]byte(secret)}
}
func (s *Signer) Sign(claims *jwt.StandardClaims) {
	return Sign(s.secret, claims)
}

func (s *Signer) Parse(token string) (*jwt.StandardClaims, error) {
	return Parse(s.secret, token)
}

func NewStandardClaims(aud, sub, iss string, iat, exp int64) *jwt.StandardClaims {
	return &jwt.StandardClaims{
		Audience:  aud, // URL of the token endpoint.
		ExpiresAt: exp,
		IssuedAt:  iat,
		Issuer:    iss, // Client ID
		Subject:   sub, // Client ID
	}
}

func Sign(key []byte, claims *jwt.StandardClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

func Parse(key []byte, token string) (*jwt.StandardClaims, error) {
	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if t == nil {
		return nil, errors.New("invalid token")
	}
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
			return nil, fmt.Errorf("token malformed: %s", err)
		}
	}
	return nil, fmt.Errorf("token malformed: %s", err.Error())
}

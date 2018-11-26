package main

import (
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

type Signer interface {
	Sign(claims jwt.Claims) (string, error)
}

type JWTSigner struct {
	secret []byte
}

func NewSigner(secret []byte) *JWTSigner {
	return &JWTSigner{secret}
}

func NewNopSigner() *JWTSigner {
	return &JWTSigner{secret: []byte("")}
}

func (s *JWTSigner) Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *JWTSigner) Parse(token string) (*jwt.StandardClaims, error) {
	return ParseJWT(s.secret, token)
}

func ParseJWT(key []byte, token string) (*jwt.StandardClaims, error) {
	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if t.Valid {
		if claims, ok := t.Claims.(*jwt.StandardClaims); ok && t.Valid {
			return claims, nil
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		switch {
		case ve.Errors&jwt.ValidationErrorMalformed != 0:
			return nil, errors.New("token is invalid")
		case ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0:
			return nil, errors.New("token expired")
		default:
			return nil, fmt.Errorf("token malformed: %s", err.Error())
		}
	}
	return nil, fmt.Errorf("token malformed: %s", err.Error())
}

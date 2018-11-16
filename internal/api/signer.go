package main

import jwt "github.com/dgrijalva/jwt-go"

type Signer interface {
	Sign(claims jwt.Claims) (string, error)
}

type JWTSigner struct {
	secret []byte
}

func NewSigner(secret []byte) *JWTSigner {
	return &JWTSigner{
		secret: secret,
	}
}
func NewNopSigner() *JWTSigner {
	return &JWTSigner{[]byte("")}
}

func (s *JWTSigner) Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

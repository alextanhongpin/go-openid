package crypto

import (
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

// var (
//         iat = time.Now().UTC()
//         exp = iat.Add(1 * time.Hour)
// )
// claims := NewStandardClaims(aud, sub, iss, iat.Unix(), exp.Unix())

func NewStandardClaims(aud, sub, iss string, iat, exp int64) *jwt.StandardClaims {
	return &jwt.StandardClaims{
		Audience:  aud, // URL of the token endpoint.
		ExpiresAt: exp,
		IssuedAt:  iat,
		Issuer:    iss, // Client ID
		Subject:   sub, // Client ID
	}
}

func NewJWT(key []byte, claims *jwt.StandardClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
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
			return nil, errors.New("invalid token")
		case ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0:
			return nil, errors.New("token expired")
		default:
			return nil, fmt.Errorf("token malformed: %s", err.Error())
		}

	}
	return nil, fmt.Errorf("token malformed: %s", err.Error())
}

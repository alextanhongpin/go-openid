package crypto

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/rs/xid"
)

// Crypto represents the encryption/decryption methods for OIDC
type Crypto interface {
	Code() string
	UUID() string
	NewJWT(aud, sub, iss string, dur time.Duration) (string, error)
	ParseJWT(token string) (*jwt.Token, error)
}

type Impl struct {
	key []byte
}

func New(key string) *Impl {
	return &Impl{
		key: []byte(key),
	}
}

func (c *Impl) Code() string {
	return xid.New().String()
}

func (c *Impl) UUID() string {
	return uuid.New().String()
}

func (c *Impl) NewJWT(aud, sub, iss string, dur time.Duration) (string, error) {
	claims := &jwt.StandardClaims{
		Audience:  aud,
		ExpiresAt: time.Now().Add(dur).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    iss,
		Subject:   sub,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(c.key)
}

func (c *Impl) ParseJWT(token string) (*jwt.Token, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return c.key, nil
	})

	if t.Valid {
		return t, nil
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

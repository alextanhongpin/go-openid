package models

import jwt "github.com/dgrijalva/jwt-go"

// Claims represents a custom jwt claim with user id
type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

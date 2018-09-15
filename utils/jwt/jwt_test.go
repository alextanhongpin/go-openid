package jwt

import (
	"fmt"
	"testing"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
)

func TestSignClaim(t *testing.T) {

	claims := Claims{
		"123456",
		jwtgo.StandardClaims{
			ExpiresAt: int64(24 * time.Hour),
			Issuer:    "www.openid",
		},
	}
	expected := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzNDU2IiwiZXhwIjo4NjQwMDAwMDAwMDAwMCwiaXNzIjoid3d3Lm9wZW5pZCJ9.nluuA2dBarCOSAKo3FDC3ZZp0LnzURk_94r55noTaoY"
	got, _ := Sign(claims)

	if expected != got {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestVerifyClaim(t *testing.T) {
	claims, err := Verify("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzNDU2IiwiZXhwIjo4NjQwMDAwMDAwMDAwMCwiaXNzIjoid3d3Lm9wZW5pZCJ9.nluuA2dBarCOSAKo3FDC3ZZp0LnzURk_94r55noTaoY")
	if err != nil {
		fmt.Println(err)
	}
	expected := int64(86400000000000)
	got := claims.StandardClaims.ExpiresAt

	if expected != got {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

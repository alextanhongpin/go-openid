package jwt

import jwtgo "github.com/dgrijalva/jwt-go"
import "fmt"

var secret = []byte("super_secret")

type Claims struct {
	UserID string `json:"user_id"`
	jwtgo.StandardClaims
}

// Verify the JWT token and return the claims
func Verify(tokenString string) (*Claims, error) {
	token, err := jwtgo.ParseWithClaims(tokenString, &Claims{}, func(token *jwtgo.Token) (interface{}, error) {
		// if _, ok := token.Method(*jwt.SigningMethodHMAC); !ok {
		// 	return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		// }
		return secret, nil
	})
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		fmt.Println(claims.UserID, claims.StandardClaims.ExpiresAt)
		return claims, nil
	}
	return &Claims{}, err
}

// Sign a token and return it
func (c Claims) Sign() (string, error) {
	// claims := Claims{
	// 	"123456",
	// 	jwt.StandardClaims{
	// 		ExpiresAt: 15000,
	// 		Issuer:    "www.openid",
	// 	},
	// }

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, c)
	ss, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return ss, nil
}

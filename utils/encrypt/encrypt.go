package encrypt

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var errPasswordRequired = errors.New("Password cannot be empty")

// HashPassword takes a string and convert it to a hash
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errPasswordRequired
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// ComparePassword compares the password and the hash
func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

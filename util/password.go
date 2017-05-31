package util

import (
  "golang.org/x/crypto/bcrypt"
)

// HashPassword takes a string and convert it to a hash
func HashPassword(password string) (string, error) {
  bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
  return string(bytes), err
}

// ComparePassword compares the password and the hash
func ComparePassword(password, hash string) bool {
  err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
  return err == nil
}

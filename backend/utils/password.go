package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	MinCost     int = 4
	MaxCost     int = 31
	DefaultCost int = 10
)

// Encrypt user password
func HashingPassword(password string) (string, error) {
	// Hashing password byte with default cost
	hashedPassByte, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	return string(hashedPassByte), err
}

// Match user password
func MatchingPassword(inputPass string, dbPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(inputPass))
	return err == nil
}

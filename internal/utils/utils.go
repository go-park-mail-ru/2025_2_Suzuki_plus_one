package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

// HashPassword hashes a plain-text password using SHA-256
func HashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", errors.New("password must be at least 8 characters long")
	}
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))
	return hashedPassword, nil
}

// ValidateHashedPassword compares a hashed password with its plain-text version
func ValidateHashedPassword(hashedPassword, password string) error {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	computedHash := hex.EncodeToString(hasher.Sum(nil))

	if hashedPassword != computedHash {
		return errors.New("password does not match")
	}
	return nil
}

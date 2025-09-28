package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// HashString hashes a string using SHA-256
func HashString(target string) (string, error) {
	hasher := sha256.New()
	hasher.Write([]byte(target))
	hashedString := hex.EncodeToString(hasher.Sum(nil))
	return hashedString, nil
}

// ValidateHashedString compares a hashed string with its plain-text version
//
// Args:
//
//	hashedString: The hashed string to compare against
//	plainText: The plain-text string to validate
func ValidateHashedString(hashedString, plainText string) error {
	hasher := sha256.New()
	hasher.Write([]byte(plainText))
	computedHash := hex.EncodeToString(hasher.Sum(nil))

	if hashedString != computedHash {
		return errors.New("string does not match")
	}
	return nil
}

func HashPasswordBcrypt(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func ValidateHashedPasswordBcrypt(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.New("password does not match")
	}
	return nil
}

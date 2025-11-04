package common

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash a password using bcrypt
// Notice that bcrypt automatically handles salting
func HashPasswordBcrypt(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// Validate a bcrypt hashed password against its plain-text version
func ValidateHashedPasswordBcrypt(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

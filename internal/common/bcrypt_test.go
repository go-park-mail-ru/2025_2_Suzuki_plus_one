package common

import "testing"

func TestHashPasswordBcrypt(t *testing.T) {
	password := "securePassword123!"
	hashedPassword, err := HashPasswordBcrypt(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	if hashedPassword == password {
		t.Fatalf("Hashed password should not be the same as the original password")
	}
}

func TestValidateHashedPasswordBcrypt(t *testing.T) {
	password := "securePassword123!"
	hashedPassword, err := HashPasswordBcrypt(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	err = ValidateHashedPasswordBcrypt(hashedPassword, password)
	if err != nil {
		t.Fatalf("Error validating hashed password: %v", err)
	}

	// Test with incorrect password
	err = ValidateHashedPasswordBcrypt(hashedPassword, "wrongPassword")
	if err == nil {
		t.Fatalf("Expected error when validating with wrong password, got nil")
	}
}

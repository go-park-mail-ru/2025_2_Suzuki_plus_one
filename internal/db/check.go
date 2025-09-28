package db

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

func purifyInputString(s string) string {
	// Remove invalid characters: < > ; ' " `
	return strings.Map(func(r rune) rune {
		switch r {
		case '<', '>', ';', '\'', '"', '`':
			return -1
		default:
			return r
		}
	}, s)
}

// Checks if the provided email is valid
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	regex := regexp.MustCompile(`^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	pureEmail := purifyInputString(email)
	parts := strings.SplitN(pureEmail, "@", 2)
	if len(parts) != 2 {
		return errors.New("incorrect email format")
	}
	localPart, domain := parts[0], parts[1]

	if pureEmail != email {
		return errors.New("email contains invalid characters (< > ; ' \" `)")
	}
	if len(pureEmail) > 256 {
		return errors.New("email is too long (maximum 256 characters)")
	}
	if !regex.MatchString(pureEmail) {
		return errors.New("incorrect email format")
	}
	for _, r := range pureEmail {
		if unicode.IsSpace(r) {
			return errors.New("email must not contain spaces")
		}
	}
	if len(localPart) > 64 {
		return errors.New("local part of the email is too long (maximum 64 characters)")
	}
	if !strings.Contains(domain, ".") {
		return errors.New("domain must contain a dot")
	}
	return nil
}

// Checks if the provided password is valid
func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}
	purePassword := purifyInputString(password)

	if purePassword != password {
		return errors.New("password contains invalid characters (< > ; ' \" `)")
	}
	if len(purePassword) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(purePassword) > 128 {
		return errors.New("password is too long (maximum 128 characters)")
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false
	specialChars := "-/=+!@#$%^&*()"
	for i, r := range purePassword {
		if unicode.IsUpper(r) {
			hasUpper = true
		}
		if unicode.IsLower(r) {
			hasLower = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
		if strings.ContainsRune(specialChars, r) {
			hasSpecial = true
		}
		if unicode.IsSpace(r) {
			return errors.New("password must not contain spaces")
		}
		// Check for more than 3 identical characters in a row
		if i >= 3 &&
			r == rune(purePassword[i-1]) &&
			r == rune(purePassword[i-2]) &&
			r == rune(purePassword[i-3]) {
			return errors.New("the password must not contain more than 3 identical characters in a row")
		}
	}
	if !hasUpper {
		return errors.New("password must contain at least one capital letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return errors.New("password must contain at least one digit")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character (/-=+!@#$%^&*())")
	}
	return nil
}

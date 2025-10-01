package utils

import (
	"strings"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email   string
		wantErr bool
	}{
		{"", true},
		{"plainaddress", true},
		{"@missinglocal.org", true},
		{"missingatsign.com", true},
		{"user@.com", true},
		{"user@domain", true},
		{"user@domain.c", true},
		{"user@domain.com", false},
		{"user.name@domain.com", false},
		{"user_name@domain.co.uk", false},
		{"user-name@domain.com", false},
		{"user+name@domain.com", true}, // + not allowed by regex
		{"user@domain..com", true},
		{"user@domain.com;", true},
		{"user@domain.com<", true},
		{"user@domain.com>", true},
		{"user@domain.com'", true},
		{"user@domain.com\"", true},
		{"user@domain.com`", true},
		{"user@domain.com ", true},
		{" user@domain.com", true},
		{"user@domain.com ", true},
		{"user@domain.com", false},
		{strings.Repeat("a", 65) + "@domain.com", true},     // local part too long
		{"user@" + strings.Repeat("a", 250) + ".com", true}, // total too long
	}

	for _, tt := range tests {
		err := ValidateEmail(tt.email)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateEmail(%q) error = %v, wantErr %v", tt.email, err, tt.wantErr)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password string
		wantErr  bool
	}{
		{"", true},
		{"short1A!", false},
		{"NoSpecialChar1", true},
		{"nouppercase1!", true},
		{"NOLOWERCASE1!", true},
		{"NoDigit!!", true},
		{"NoSpace 1!", true},
		{"ValidPass1!", false},
		{"Va!1" + strings.Repeat("a", 124), true}, // more than 3 identical
		{"Va!1" + strings.Repeat("a", 130), true}, // too long
		{"aaaaAAAA1!", true},                      // 4 identical in a row
		{"Abcdefg1!", false},
		{"Abcdefg1-", false},
		{"Abcdefg1/", false},
		{"Abcdefg1=", false},
		{"Abcdefg1+", false},
		{"Abcdefg1!", false},
		{"Abcdefg1@", false},
		{"Abcdefg1#", false},
		{"Abcdefg1$", false},
		{"Abcdefg1%", false},
		{"Abcdefg1^", false},
		{"Abcdefg1&", false},
		{"Abcdefg1*", false},
		{"Abcdefg1(", false},
		{"Abcdefg1)", false},
		{"Abcdefg1;", true},
		{"Abcdefg1<", true},
		{"Abcdefg1>", true},
		{"Abcdefg1'", true},
		{"Abcdefg1\"", true},
		{"Abcdefg1`", true},
	}

	for _, tt := range tests {
		err := ValidatePassword(tt.password)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidatePassword(%q) error = %v, wantErr %v", tt.password, err, tt.wantErr)
		}
	}
}

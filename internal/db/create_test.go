package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	db := NewMockDB()

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{
			name:     "Email_already_exists",
			email:    "test@example.com",
			password: "Pa!ssword456",
			wantErr:  true,
		},
		{
			name:     "Password_too_short",
			email:    "test@example.com",
			password: "short",
			wantErr:  true,
		},
		{
			name:     "Invalid_email_format",
			email:    "invalid-email",
			password: "Pa!ssword456",
			wantErr:  true,
		},
		{
			name:     "Blank_input",
			email:    "",
			password: "",
			wantErr:  true,
		},
		{
			name:     "Valid_input",
			email:    "REALLYREALLY@example.com",
			password: "Pa!ssword456",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			user, err, created := db.CreateUser(tt.email, tt.password)
			if tt.wantErr {
				require.NotNil(t, err, "expected error but got none")
				require.False(t, created, "expected user not to be created")
			} else {
				assert.NotNil(t, user, "expected user to be non-nil")
				require.Nil(t, err, "expected no error but got one")
				require.True(t, created, "expected user to be created")

			}
		})
	}
}

package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJWT(t *testing.T) {
	InitJWT("secret")
	token, err := GenerateToken(1, AccessTokenTTL)
	require.NoError(t, err, "Failed to generate token")

	userID, err := ValidateToken(token)
	require.NoError(t, err, "Failed to validate token: %v", err)
	require.Equal(t, uint(1), userID, "UserID does not match")
}

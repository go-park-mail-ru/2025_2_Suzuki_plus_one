package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

const SECRET = "secret"
const TOKEN = "test_token"

func TestTokenHeader(t *testing.T) {
	logger := zaptest.NewLogger(t)
	defer logger.Sync()

	authenticator := NewAuthHeader(SECRET, logger)
	// Convert to TokenHeader to access its fields
	th, ok := authenticator.TokenMgr.(*TokenHeader)
	require.True(t, ok)

	// Test Exists
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	require.False(t, authenticator.TokenMgr.Exists(req), "Expected Exists to return false")

	// Inject token into header
	req.Header.Set(th.HeaderName, th.Prefix+TOKEN)
	require.True(t, authenticator.TokenMgr.Exists(req), "Expected Exists to return true")

	// Test Get
	retrievedToken := authenticator.TokenMgr.Get(req)
	require.Equal(t, TOKEN, retrievedToken, "Expected Get to return correct token")

	// Test Set
	newToken := "new_test_token"
	th.Set(req, newToken)
	retrievedToken = th.Get(req)
	require.Equal(t, newToken, retrievedToken, "Expected Get to return new token after Set")
}

func TestTokenCookie(t *testing.T) {
	logger := zaptest.NewLogger(t)
	defer logger.Sync()

	authenticator := NewAuthCookie(SECRET, logger)
	tc, ok := authenticator.TokenMgr.(*TokenCookie)
	require.True(t, ok)
	cookieName := tc.CookieName

	// Test Exists
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	require.False(t, tc.Exists(req), "Expected Exists to return false")

	// Inject token into cookies
	req.AddCookie(&http.Cookie{
		Name:  cookieName,
		Value: TOKEN,
	})
	require.True(t, tc.Exists(req), "Expected Exists to return true")

	// Test Get
	retrievedToken := tc.Get(req)
	require.Equal(t, TOKEN, retrievedToken, "Expected Get to return correct token")

	// Test Set
	newToken := "new_test_token"
	tc.Set(req, newToken)
	retrievedToken = tc.Get(req)
	require.Equal(t, newToken, retrievedToken, "Expected Get to return new token after Set")
}

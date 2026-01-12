package handlers

import (
	"net/http"
	//"net/http"
	"testing"
	"time"

	//"time"

	"github.com/stretchr/testify/assert"
)

func TestCookieConstants(t *testing.T) {
	assert.Equal(t, "refresh_token", CookieRefreshTokenName)
}

func TestNewResetCookieRefreshToken(t *testing.T) {
	cookie := NewResetCookieRefreshToken()

	assert.Equal(t, CookieRefreshTokenName, cookie.Name)
	assert.Equal(t, "", cookie.Value)
	assert.Equal(t, "/", cookie.Path)
	assert.True(t, cookie.HttpOnly)
	assert.False(t, cookie.Secure) // TODO: change in production
	assert.Equal(t, time.Unix(0, 0), cookie.Expires)
	assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
	assert.True(t, cookie.Expires.Before(time.Now()), "Cookie should be expired")
}

func TestNewRefreshTokenCookie(t *testing.T) {
	testToken := "test-refresh-token-123"
	cookie := NewRefreshTokenCookie(testToken)

	assert.Equal(t, CookieRefreshTokenName, cookie.Name)
	assert.Equal(t, testToken, cookie.Value)
	assert.Equal(t, "/", cookie.Path)
	assert.True(t, cookie.HttpOnly)
	assert.False(t, cookie.Secure) // TODO: change in production
	assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)

	expectedExpiry := time.Now().Add(7 * 24 * time.Hour)
	timeDiff := expectedExpiry.Sub(cookie.Expires)
	assert.True(t, timeDiff < time.Second, "Cookie expiry should be approximately 7 days from now")
	assert.True(t, cookie.Expires.After(time.Now()), "Cookie should not be expired")
}

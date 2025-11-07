package handlers

import (
	"net/http"
	"time"
)

var CookieRefreshTokenName = "refresh_token"

func NewResetCookieRefreshToken() *http.Cookie {
	return &http.Cookie{
		Name:     CookieRefreshTokenName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // TODO: set to true in production
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
	}
}

func NewRefreshTokenCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     CookieRefreshTokenName,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,                              // TODO: set to true in production
		Expires:  time.Now().Add(7 * 24 * time.Hour), // TODO: change from 7 days if needed
		SameSite: http.SameSiteStrictMode,
	}
}

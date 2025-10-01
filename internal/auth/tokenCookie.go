package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/models"
)

type TokenCookie struct {
	CookieName string // Name of the cookie, e.g. "refresh_token"
}

// Returns new TokenCookie with default cookie name "token"
func NewTokenCookie() *TokenCookie {
	return &TokenCookie{
		CookieName: "TOKEN",
	}
}

// Checks if the cookie exists in the request
// Returns true if exists, false otherwise
func (tc *TokenCookie) Exists(request *http.Request) bool {
	_, err := request.Cookie(tc.CookieName)
	return err == nil
}

// Retrieves token string from the cookie
// Returns empty string if no cookie found
func (tc *TokenCookie) Get(request *http.Request) string {
	cookie, err := request.Cookie(tc.CookieName)
	log.Printf("Get cookie:[%s]=[%v] err=[%v]", tc.CookieName, cookie, err)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// ...existing code...
func (tc *TokenCookie) Set(request *http.Request, token string) {
	// default request.AddCookie does not replace existing cookies,

	// Remove any existing cookie with the same name
	cookies := request.Header.Get("Cookie")
	newCookie := fmt.Sprintf("%s=%s", tc.CookieName, token)
	var filtered []string
	for c := range strings.SplitSeq(cookies, ";") {
		c = strings.TrimSpace(c)
		if !strings.HasPrefix(c, tc.CookieName+"=") {
			filtered = append(filtered, c)
		}
	}
	filtered = append(filtered, newCookie)
	request.Header.Set("Cookie", strings.Join(filtered, "; "))
}

// Sets token in the cookie and sends response as JSON
func (tc *TokenCookie) ResponseWithAuth(writer http.ResponseWriter, token string, response models.SignInResponse) {
	// Set token in the cookie
	cookie := tc.GetCookieWithName(token)
	// Set cookie expiration
	cookie.Expires = time.Now().Add(24 * time.Hour)
	http.SetCookie(writer, cookie)
	log.Println("auth: Cookie set in ResponseWithAuth")
	json.NewEncoder(writer).Encode(response)
}

// Send clear cookie
func (tc *TokenCookie) ResponseWithDeauth(writer http.ResponseWriter) {
	// Clear the cookie by setting its expiration date to the past
	cookie := tc.GetCookieWithName("")
	cookie.MaxAge = -1
	log.Println("auth: Cookie cleared in ResponseWithDeauth")
	http.SetCookie(writer, cookie)
}

// TODO: improve cookie settings
func (tc *TokenCookie) GetCookieWithName(token string) *http.Cookie {
	return &http.Cookie{
		Name:     tc.CookieName,
		Value:    token,
		Path:     "/",                  // Ensure cookie is sent for all paths
		HttpOnly: true,                 // Prevent JS access to cookie
		Secure:   false,                // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode, // Helps with cross-site requests
	}
}

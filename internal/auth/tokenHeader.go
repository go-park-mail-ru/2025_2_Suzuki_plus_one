package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/models"
)

type TokenHeader struct {
	HeaderName string // Field name in header, e.g. "Authorization"
	Prefix     string // Prefix before token, e.g. "Bearer "
}

func NewTokenHeader() *TokenHeader {
	return &TokenHeader{
		HeaderName: "Authorization",
		Prefix:     "Bearer ",
	}
}

// Checks if "Authorization" header exists in request
// Returns true if exists, false otherwise
func (th *TokenHeader) Exists(request *http.Request) bool {
	authHeader := request.Header.Get(th.HeaderName)
	return authHeader != ""
}

// Retrieves token string from "Authorization" header
// Expects header in format "Bearer <token>"
//
// Returns empty string if no token found or token has a wrong format
func (th *TokenHeader) Get(request *http.Request) string {
	authHeader := request.Header.Get(th.HeaderName)
	// Trim "Bearer " prefix
	if len(authHeader) > len(th.Prefix) &&
		authHeader[:len(th.Prefix)] == th.Prefix {
		return authHeader[len(th.Prefix):]
	}
	return ""
}

// Sets token string in "Authorization" header (adds "Bearer " prefix)
func (th *TokenHeader) Set(request *http.Request, token string) {
	request.Header.Set(th.HeaderName, th.Prefix+token)
}

// Set token in the token field and sends response as JSON
func (th *TokenHeader) ResponseWithAuth(writer http.ResponseWriter, token string, response models.SignInResponse) {
	// Set token in the header
	// writer.Header().Set("Access-Control-Expose-Headers", th.HeaderName)
	// writer.Header().Set(th.HeaderName, th.Prefix+token)

	// TODO: this is disabled and propably should be removed
	// But it can be useful for access tokens

	// Set token in response body
	response_with_token := struct {
		models.SignInResponse
		Token string `json:"token"`
	}{
		SignInResponse: response,
		Token:          token,
	}
	json.NewEncoder(writer).Encode(response_with_token)
}

// ! Useless function, because we cannot really log out user by just removing the header
func (th *TokenHeader) ResponseWithDeauth(writer http.ResponseWriter) {
	// Note: this does not really log out the user on client side,
}

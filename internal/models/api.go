// API request and response models
// Which have to correspond to OpenAPI spec

// Note: Probably we can add http method as a suffix to struct names

package models

import "net/http"

// # Requests
type SignUpRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Username  string `json:"username"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MoviesRequest struct {
	Limit  uint `query:"limit,omitempty"`
	Offset uint `query:"offset,omitempty"`
}

// # Responses

type Responder interface {
	Respond(w http.ResponseWriter) error
}

type SignInResponse struct {
	User UserAPI `json:"user"`
}

type SignUpResponse struct {
	User UserAPI `json:"user"`
}

type AuthResponse struct {
	User UserAPI `json:"user"`
}

// ## Error
type ErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message,omitempty"`
	Details string `json:"details,omitempty"` // Internal error details, optional
}

// API request and response models
// Which have to correspond to OpenAPI spec

package models

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

// # Responses
type SignInResponse struct {
	Token string  `json:"token"`
	User  UserAPI `json:"user"`
}

// ## Error
type ErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message,omitempty"`
	Details string `json:"details,omitempty"` // Internal error details, optional
}

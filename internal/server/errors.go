package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/models"
)

var (
	ErrNoTokenProvided  = models.ErrorResponse{Type: "auth", Message: "no token provided"}
	ErrInvalidOrExpired = models.ErrorResponse{Type: "auth", Message: "invalid or expired token"}
	ErrSignInWrongData  = models.ErrorResponse{Type: "auth", Message: "wrong email or password"}
	ErrSignInInternal   = models.ErrorResponse{Type: "auth", Message: "internal server error"}
	ErrSignUpWrongData  = models.ErrorResponse{Type: "auth", Message: "wrong email or password"}
	ErrSignUpUserExists = models.ErrorResponse{Type: "auth", Message: "user does not exist"}
	ErrSignUpInternal   = models.ErrorResponse{Type: "auth", Message: "internal server error"}

	ErrMoviesInvalidParams = models.ErrorResponse{Type: "content", Message: "invalid request parameters"}
)

// Helper function to respond with an error in JSON format
func responseWithError(w http.ResponseWriter, statusCode int, err models.ErrorResponse) {
	log.Printf("Error: %s - %s", err.Type, err.Message)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(err)
}

// Wrapper to add details to an error response and send it with [server.responseWithError]
func errorWithDetails(err models.ErrorResponse, details string) models.ErrorResponse {
	err.Details = details
	return err
}

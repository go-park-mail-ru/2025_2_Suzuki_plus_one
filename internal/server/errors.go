package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/models"
	"go.uber.org/zap"
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
	ErrActorIdIsRequired   = models.ErrorResponse{Type: "content", Message: "ActorID is required"}
	ErrACtorNotFound       = models.ErrorResponse{Type: "content", Message: "Actor not found"}
)

// Helper function to respond with an error in JSON format
func responseWithError(w http.ResponseWriter, statusCode int, err models.ErrorResponse, logger *zap.Logger) {
	logger.Error("HTTP error response",
		zap.Int("status_code", statusCode),
		zap.String("error_type", err.Type),
		zap.String("error_message", err.Message),
	)

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(err)
}

// Wrapper to add details to an error response and send it with [server.responseWithError]
func errorWithDetails(err models.ErrorResponse, details string) models.ErrorResponse {
	err.Details = details
	return err
}

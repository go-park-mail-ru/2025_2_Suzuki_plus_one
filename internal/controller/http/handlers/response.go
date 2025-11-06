package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type Response struct {
	Code int
	Data any
}

// [dto.Error] wrapper for HTTP handlers
type ResponseError struct {
	Code    int
	Message error
	Details string
}

// Encodes and sends a JSON response with the given status code and data
func Respond(l logger.Logger, w http.ResponseWriter, status int, data any) {
	l.Info("HTTP response",
		l.ToInt("status_code", status),
	)

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Sends a JSON error ResponseError with override details
func RespondWithError(l logger.Logger, w http.ResponseWriter, err ResponseError, details string) {
	dto := dto.Error{
		Type:    "controller/http",
		Message: err.Message.Error(),
		Details: details,
	}

	l.Error("Error HTTP response",
		l.ToInt("status_code", err.Code),
		l.ToString("error_type", dto.Type),
		l.ToString("error_message", dto.Message),
		l.ToString("error_details", dto.Details),
	)

	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(dto)
}

// Sends a JSON error response using a dto.Error
func RespondWithDTOError(l logger.Logger, w http.ResponseWriter, err ResponseError, dtoErr *dto.Error) {
	l.Error("Error HTTP response",
		l.ToInt("status_code", err.Code),
		l.ToString("error_type", dtoErr.Type),
		l.ToString("error_message", dtoErr.Message),
		l.ToString("error_details", dtoErr.Details),
	)

	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(dtoErr)
}

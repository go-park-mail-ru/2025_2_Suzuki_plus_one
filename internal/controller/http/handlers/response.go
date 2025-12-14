package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/mailru/easyjson"
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
	if err := json.NewEncoder(w).Encode(data); err != nil {
		l.Error("Failed to encode response", l.ToError(err))
	}
}

// Sends a JSON error ResponseError with overridden details
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
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		l.Error("Failed to encode error response", l.ToError(err))
	}
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
	if err := json.NewEncoder(w).Encode(dtoErr); err != nil {
		l.Error("Failed to encode dto error response", l.ToError(err))
	}
}

// Easy json
func RespondEasyJSON(l logger.Logger, w http.ResponseWriter, status int, data easyjson.Marshaler) {
	l.Info("HTTP response",
		l.ToInt("status_code", status),
	)

	// Set status code (will be overwritten to 500 on write error)
	w.WriteHeader(status)

	// Serialize
	_, _, err := easyjson.MarshalToHTTPResponseWriter(
		data,
		w,
	)
	// Flush status code and data
	if err != nil {
		l.Error("Failed to encode response", l.ToError(err))
	}
}

func RespondWithDTOErrorEasyJSON(l logger.Logger, w http.ResponseWriter, err ResponseError, dtoErr easyjson.Marshaler) {
	// Set status code (will be overwritten to 500 on write error)
	w.WriteHeader(err.Code)

	// Serialize
	_, _, mErr := easyjson.MarshalToHTTPResponseWriter(
		dtoErr,
		w,
	)
	// Flush status code and data
	if mErr != nil {
		l.Error("Failed to encode dto error response", l.ToError(mErr))
	}
}

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
)

// [dto.Error] wrapper for HTTP handlers
type ResponseError struct {
	Code    int
	Message string
}

// Encodes and sends a JSON response with the given status code and data
func (h *Handlers) Response(w http.ResponseWriter, status int, data any) {
	h.Logger.Info("HTTP response",
		h.Logger.ToInt("status_code", status),
	)

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Sets error type and sends a JSON error response
func (h *Handlers) ResponseWithError(w http.ResponseWriter, err ResponseError, details string) {
	dto := dto.Error{
		Type:    "controller/http",
		Message: err.Message,
		Details: details,
	}

	h.Logger.Error("HTTP response",
		h.Logger.ToInt("status_code", err.Code),
		h.Logger.ToString("error_type", dto.Type),
		h.Logger.ToString("error_message", dto.Message),
		h.Logger.ToString("error_details", dto.Details),
	)

	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(dto)
}

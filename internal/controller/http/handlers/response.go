package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
)

// Encodes and sends a JSON response with the given status code and data
func (h *Handlers) Response(w http.ResponseWriter, status int, data any) {
	h.logger.Info("HTTP response",
		h.logger.ToInt("status_code", status),
	)

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handlers) ResponseWithError(w http.ResponseWriter, err ResponseError, details string) {
	dto := dto.Error{
		Type:    "controller/http",
		Message: err.Message,
		Details: details,
	}

	h.logger.Error("HTTP response",
		h.logger.ToInt("status_code", err.Code),
		h.logger.ToString("error_type", dto.Type),
		h.logger.ToString("error_message", dto.Message),
		h.logger.ToString("error_details", dto.Details),
	)

	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(dto)
}

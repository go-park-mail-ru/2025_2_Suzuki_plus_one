package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
)

var (
	ErrGetMediaWatchInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: "Invalid parameters for media fetching",
	}
)

var URLParamMediaWatchID = "media_id"

// Get all movies from database
func (h *Handlers) GetMediaWatch(w http.ResponseWriter, r *http.Request) {
	// Handle input parameters
	input := dto.GetMediaWatchInput{}
	rp := NewRequestParams(h.Logger, r, &input)
	rp.AddQuery(URLParamMediaWatchID, &input.MediaID)

	if err := rp.Parse(); err != nil {
		h.Logger.Error("Failed to parse query parameters",
			h.Logger.ToString("error", err.Error()))
		h.ResponseWithError(w, ErrGetMediaWatchInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.GetMediaWatchUseCase.Execute(rp.GetContext(), input)
	if err != nil {
		h.Logger.Error("Failed to fetch media",
			h.Logger.ToString("error", err.Message))
		// Respond with error
		h.Response(w, ErrGetMediaWatchInvalidParams.Code, err)
		return
	}

	h.Logger.Debug("Fetching media link completed successfully",
		h.Logger.ToString("media_id", strconv.FormatUint(uint64(input.MediaID), 10)),
		h.Logger.ToString("link", output.URL))

	// Respond with output
	h.Response(w, http.StatusOK, output)
}

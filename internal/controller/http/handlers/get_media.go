package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
)

var (
	ErrGetMediaInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: "Invalid parameters for media fetching",
	}
)

const URLParamMediaID = "media_id"

// Get all movies from database
func (h *Handlers) GetMedia(w http.ResponseWriter, r *http.Request) {
	// Handle input parameters
	input := dto.GetMediaInput{}
	rp := NewRequestParams(h.Logger, r, &input)
	rp.AddPath(URLParamMediaID, &input.MediaID)

	if err := rp.Parse(); err != nil {
		h.Logger.Error("Failed to parse query parameters",
			h.Logger.ToString("error", err.Error()))
		h.ResponseWithError(w, ErrGetMediaInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.GetMediaUseCase.Execute(rp.GetContext(), input)
	if err != nil {
		h.Logger.Error("Failed to fetch media",
			h.Logger.ToString("error", err.Message))
		// Respond with error
		h.Response(w, ErrGetMediaInvalidParams.Code, err)
		return
	}

	h.Logger.Debug("Fetching media details completed successfully",
		h.Logger.ToString("media_id", strconv.FormatUint(uint64(input.MediaID), 10)),
		h.Logger.ToString("actor_count", strconv.Itoa(len(output.Actors))))

	// Respond with output
	h.Response(w, http.StatusOK, output)
}

package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
)

var (
	ErrGetActorInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: "Invalid parameters for actor",
	}
)

const URLParamActorID = "actor_id"

// Get all movies from database
func (h *Handlers) GetActor(w http.ResponseWriter, r *http.Request) {
	// Handle input parameters
	input := dto.GetActorInput{}
	rp := NewRequestParams(h.Logger, r, &input)
	rp.AddPath(URLParamActorID, &input.ActorID)

	if err := rp.Parse(); err != nil {
		h.Logger.Error("Failed to parse query parameters",
			h.Logger.ToString("error", err.Error()))
		h.ResponseWithError(w, ErrGetActorInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.GetActorUseCase.Execute(rp.GetContext(), input)
	if err != nil {
		h.Logger.Error("Failed to fetch actor",
			h.Logger.ToString("error", err.Message))
		// Respond with error
		h.Response(w, ErrMoviesInvalidParams.Code, err)
		return
	}

	h.Logger.Debug("Fetching actor details completed successfully",
		h.Logger.ToString("actor_id", strconv.FormatUint(uint64(input.ActorID), 10)),
		h.Logger.ToString("media_count", strconv.Itoa(len(output.Medias))))

	// Respond with output
	h.Response(w, http.StatusOK, output)
}

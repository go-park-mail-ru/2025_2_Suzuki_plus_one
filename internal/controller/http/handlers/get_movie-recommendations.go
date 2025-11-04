package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
)

var (
	ErrMoviesInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: "Invalid parameters for movies",
	}
)

// Get all movies from database
func (h *Handlers) GetMovieRecommendations(w http.ResponseWriter, r *http.Request) {
	// Handle input parameters
	input := dto.GetMovieRecommendationsInput{}
	rp := NewRequestParams(h.Logger, r, &input)
	rp.AddQuery("offset", &input.Offset)
	rp.AddQuery("limit", &input.Limit)
	if err := rp.Parse(); err != nil {
		h.Logger.Error("Failed to parse query parameters",
			h.Logger.ToString("error", err.Error()))
		h.ResponseWithError(w, ErrMoviesInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.GetMovieRecommendationsUseCase.Execute(rp.GetContext(), input)
	if err != nil {
		h.Logger.Error("Failed to fetch movie recommendations",
			h.Logger.ToString("error", err.Message))
		// Respond with error
		h.Response(w, ErrMoviesInvalidParams.Code, err)
		return
	}

	h.Logger.Info("Fetching movie recommendations completed successfully",
		h.Logger.ToString("count", strconv.FormatInt(int64(len(output.Movies)), 10)),
		h.Logger.ToString("offset", strconv.FormatUint(uint64(input.Offset), 10)),
		h.Logger.ToString("limit", strconv.FormatUint(uint64(input.Limit), 10)))

	// Respond with output
	h.Response(w, http.StatusOK, output)
}

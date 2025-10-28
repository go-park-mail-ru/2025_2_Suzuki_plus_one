package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
)

var (
	ErrMoviesInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: "Invalid query parameters for movies",
	}
)

// Get all movies from database
func (h *Handlers) GetMovies(w http.ResponseWriter, r *http.Request) {
	input := dto.GetMoviesInput{}

	// TODO: add validation middleware
	// Parse query parameters
	query := r.URL.Query()
	if offStr := query.Get("offset"); offStr != "" {
		// If parameter is blank we leave it as default 0
		if _, err := fmt.Sscanf(offStr, "%d", &input.Offset); err != nil {
			h.logger.Warn("Invalid offset parameter",
				h.logger.ToString("offset", offStr),
				h.logger.ToError(err))

			h.ResponseWithError(w, ErrMoviesInvalidParams, "Invalid offset parameter")
			return
		}
	}
	if limStr := query.Get("limit"); limStr != "" {
		// If parameter is blank we leave it as default 0 (means no limit)
		if _, err := fmt.Sscanf(limStr, "%d", &input.Limit); err != nil {
			h.logger.Warn("Invalid limit parameter",
				h.logger.ToString("limit", limStr),
				h.logger.ToError(err))
			h.ResponseWithError(w, ErrMoviesInvalidParams, "Invalid limit parameter")
			return
		}
	}

	// Execute use case
	output, err := h.GetMoviesUseCase.Execute(input)
	if err != nil {
		h.logger.Error("Failed to fetch movies",
			h.logger.ToString("error", err.Message))
		h.ResponseWithError(w, ErrMoviesInvalidParams, "Failed to fetch movies")
		return
	}

	h.logger.Info("Fetching completed successfully",
		h.logger.ToString("count", strconv.FormatInt(int64(len(output.Movies)), 10)),
		h.logger.ToString("offset", strconv.FormatUint(uint64(input.Offset), 10)),
		h.logger.ToString("limit", strconv.FormatUint(uint64(input.Limit), 10)))

	h.Response(w, http.StatusOK, output)
}

package handlers

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for GetGenreAll handler
var (
	ErrGetGenreAllInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("Invalid parameters for genre all fetching"),
	}
	ResponseGetGenreAll = Response{
		Code: http.StatusOK,
	}
)

// GetGenreAll handler
func (h *Handlers) GetGenreAll(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.GetGenreAllInput{}
	rp := NewRequestParams(log, r, &input)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse query parameters",
			log.ToString("error", err.Error()),
		)
		RespondWithError(log, w, ErrGetGenreAllInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.GetGenreAllUseCase.Execute(ctx, input)
	if err != nil {
		log.Error(
			"Failed to fetch genre",
			log.ToString("error", err.Message),
		)
		RespondWithDTOError(log, w, ErrGetGenreAllInvalidParams, err)
		return
	}

	log.Debug(
		"GetGenreAll succeeded",
		log.ToInt("genres_count", len(output.Genres)),
	)

	// Respond with output
	Respond(log, w, ResponseGetGenreAll.Code, output)
}

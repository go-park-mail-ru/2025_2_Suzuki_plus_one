package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for GetGenre handler
var (
	ErrGetGenreInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("invalid parameters for genre fetching"),
	}
	ResponseGetGenre = Response{
		Code: http.StatusOK,
	}
)

// Path parameters
const (
	PathParamGetGenreID = "genre_id"
)

// GetGenre handler
func (h *Handlers) GetGenre(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.GetGenreInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddPath(PathParamGetGenreID, &input.GenreID)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse query parameters",
			log.ToString("error", err.Error()),
		)
		RespondWithError(log, w, ErrGetGenreInvalidParams, err.Error())
		return
	}
	log.Debug(
		"GetGenre called for ID",
		log.ToString(PathParamGetGenreID, strconv.FormatUint(uint64(input.GenreID), 10)),
	)

	// Execute use case
	output, err := h.GetGenreUseCase.Execute(ctx, input)
	if err != nil {
		log.Error(
			"Failed to fetch genre",
			log.ToString("error", err.Message),
		)
		RespondWithDTOError(log, w, ErrGetGenreInvalidParams, err)
		return
	}

	log.Debug(
		"GetGenre succeeded",
		log.ToString(PathParamGetGenreID, strconv.FormatUint(uint64(input.GenreID), 10)),
	)

	// Respond with output
	Respond(log, w, ResponseGetGenre.Code, output)
}

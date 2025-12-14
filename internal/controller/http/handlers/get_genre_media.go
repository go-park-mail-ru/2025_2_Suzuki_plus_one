package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for GetGenreMedia handler
// Blank fields are not used and are filled in the handler
var (
	ErrGetGenreMediaInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("invalid parameters for genre"),
	}
	ErrGetGenreMediaNotFound = ResponseError{
		Code: http.StatusBadRequest,
	}
	ResponseGetGenreMedia = Response{
		Code: http.StatusOK,
	}
)

// Input path parameter
const PathParamGetGenreMediaID = "genre_id"

// Query parameters
const QueryParamGetGenreMediaLimit = "limit"
const QueryParamGetGenreMediaOffset = "offset"

// Get all movies from database
func (h *Handlers) GetGenreMedia(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.GetGenreMediaInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddPath(PathParamGetGenreMediaID, &input.GenreID)
	rp.AddQuery(QueryParamGetGenreMediaLimit, &input.Limit)
	rp.AddQuery(QueryParamGetGenreMediaOffset, &input.Offset)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse query parameters",
			log.ToString("error", err.Error()),
		)
		// Respond with error, if input parameters are invalid
		RespondWithError(log, w, ErrGetGenreMediaInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.GetGenreMediaUseCase.Execute(ctx, input)
	if err != nil {
		log.Error("Failed to fetch genre",
			log.ToString("error", err.Message),
		)
		// Respond with error, if use case execution fails
		RespondWithDTOError(log, w, ErrGetGenreMediaInvalidParams, err)
		return
	}

	// Log successful completion
	log.Debug(
		"Fetching genre medias completed successfully",
		log.ToString("genre_id", strconv.FormatUint(uint64(input.GenreID), 10)),
	)

	// Respond with output
	Respond(log, w, ResponseGetGenreMedia.Code, output)
}

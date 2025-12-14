package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for GetMediaEpisodes handler
// Blank fields are not used and are filled in the handler
var (
	ErrGetMediaEpisodesInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("invalid parameters for episodes"),
	}
	ErrGetMediaEpisodesNotFound = ResponseError{
		Code: http.StatusBadRequest,
	}
	ResponseGetMediaEpisodes = Response{
		Code: http.StatusOK,
	}
)

// Input path parameter
var (
	PathParamGetMediaEpisodesID = "media_id"
)

// Get all movies from database
func (h *Handlers) GetMediaEpisodes(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.GetMediaEpisodesInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddPath(PathParamGetMediaEpisodesID, &input.MediaID)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse query parameters",
			log.ToString("error", err.Error()),
		)
		// Respond with error, if input parameters are invalid
		RespondWithError(log, w, ErrGetMediaEpisodesInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.GetMediaEpisodesUseCase.Execute(ctx, input)
	if err != nil {
		log.Error("Failed to fetch episodes",
			log.ToString("error", err.Message),
		)
		// Respond with error, if use case execution fails
		RespondWithDTOError(log, w, ErrGetMediaEpisodesInvalidParams, err)
		return
	}

	// Log successful completion
	log.Debug(
		"Handler completed successfully",
		log.ToInt("media_id", int(input.MediaID)),
		log.ToString("episodes_count", strconv.Itoa(len(output.Episodes))),
	)

	// Respond with output
	Respond(log, w, ResponseGetMediaEpisodes.Code, output)
}

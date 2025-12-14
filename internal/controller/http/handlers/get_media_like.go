package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for GetMediaLike handler
// Blank fields are not used and are filled in the handler
var (
	ErrGetMediaLikeInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("invalid parameters for like"),
	}
	ErrGetMediaLikeNotFound = ResponseError{
		Code: http.StatusBadRequest,
	}
	ResponseGetMediaLike = Response{
		Code: http.StatusOK,
	}
)

// Input path parameter
const PathParamGetMediaLikeID = "media_id"

// Get all movies from database
func (h *Handlers) GetMediaLike(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.GetMediaLikeInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddPath(PathParamGetMediaLikeID, &input.MediaID)
	rp.AddAuthHeader(&input.AccessToken)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse query parameters",
			log.ToString("error", err.Error()),
		)
		// Respond with error, if input parameters are invalid
		RespondWithError(log, w, ErrGetMediaLikeInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.GetMediaLikeUseCase.Execute(ctx, input)
	if err != nil {
		log.Error("Failed to fetch actor",
			log.ToString("error", err.Message),
		)
		// Respond with error, if use case execution fails
		RespondWithDTOError(log, w, ErrGetMediaLikeInvalidParams, err)
		return
	}

	// Log successful completion
	log.Debug(
		"Handler completed successfully",
		log.ToString("media_id", strconv.Itoa(int(input.MediaID))),
		log.ToAny("liked", output.Liked),
		log.ToAny("is_dislike", output.IsDislike),
	)

	// Respond with output
	Respond(log, w, ResponseGetMediaLike.Code, output)
}

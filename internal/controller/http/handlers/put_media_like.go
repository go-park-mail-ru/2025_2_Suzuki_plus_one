package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for PutMediaLike handler
// Blank fields are not used and are filled in the handler
var (
	ErrPutMediaLikeInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("invalid parameters for like"),
	}
	ErrPutMediaLikeNotFound = ResponseError{
		Code: http.StatusBadRequest,
	}
	ResponsePutMediaLike = Response{
		Code: http.StatusOK,
	}
)

// Input path parameter
const PathParamPutMediaLikeID = "media_id"
const PathParamPutMediaLikeIsDislike = "is_dislike"

// Get all movies from database
func (h *Handlers) PutMediaLike(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.PutMediaLikeInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddPath(PathParamPutMediaLikeID, &input.MediaID)
	rp.AddPath(PathParamPutMediaLikeIsDislike, &input.IsDislike)
	rp.AddAuthHeader(&input.AccessToken)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse query parameters",
			log.ToString("error", err.Error()),
		)
		// Respond with error, if input parameters are invalid
		RespondWithError(log, w, ErrPutMediaLikeInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.PutMediaLikeUseCase.Execute(ctx, input)
	if err != nil {
		log.Error("Failed to fetch actor",
			log.ToString("error", err.Message),
		)
		// Respond with error, if use case execution fails
		RespondWithDTOError(log, w, ErrPutMediaLikeInvalidParams, err)
		return
	}

	// Log successful completion
	log.Debug(
		"Handler completed successfully",
		log.ToAny("media_id", input.MediaID),
		log.ToString("liked", strconv.FormatBool(output.Liked)),
	)

	// Respond with output
	Respond(log, w, ResponsePutMediaLike.Code, output)
}

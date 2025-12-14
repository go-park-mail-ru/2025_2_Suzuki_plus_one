package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for GetMediaMy handler
// Blank fields are not used and are filled in the handler
var (
	ErrGetMediaMyInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("invalid parameters for like"),
	}
	ErrGetMediaMyNotFound = ResponseError{
		Code: http.StatusBadRequest,
	}
	ResponseGetMediaMy = Response{
		Code: http.StatusOK,
	}
)

// Input path parameter
var (
	QueryParamGetMediaMyIsDislike = "is_dislike"
	QueryParamGetMediaMyLimit     = "limit"
	QueryParamGetMediaMyOffset    = "offset"
)

// Get all movies from database
func (h *Handlers) GetMediaMy(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.GetMediaMyInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddQuery(QueryParamGetMediaMyIsDislike, &input.IsDislike)
	rp.AddQuery(QueryParamGetMediaMyLimit, &input.Limit)
	rp.AddQuery(QueryParamGetMediaMyOffset, &input.Offset)
	rp.AddAuthHeader(&input.AccessToken)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse query parameters",
			log.ToString("error", err.Error()),
		)
		// Respond with error, if input parameters are invalid
		RespondWithError(log, w, ErrGetMediaMyInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.GetMediaMyUseCase.Execute(ctx, input)
	if err != nil {
		log.Error("Failed to fetch actor",
			log.ToString("error", err.Message),
		)
		// Respond with error, if use case execution fails
		RespondWithDTOError(log, w, ErrGetMediaMyInvalidParams, err)
		return
	}

	// Log successful completion
	log.Debug(
		"Handler completed successfully",
		log.ToString("is_dislike", strconv.FormatBool(input.IsDislike)),
		log.ToString("medias_count", strconv.Itoa(len(output.Medias))),
	)

	// Respond with output
	Respond(log, w, ResponseGetMediaMy.Code, output)
}

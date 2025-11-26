package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for GetMediaWatch handler
var (
	ErrGetMediaWatchInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("Invalid parameters for media fetching"),
	}
	ResponseGetMediaWatch = Response{
		Code: http.StatusOK,
	}
)

// URL parameter context name
const QueryParamMediaWatchID = "media_id"

// GetMediaWatch handler
func (h *Handlers) GetMediaWatch(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.GetMediaWatchInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddQuery(QueryParamMediaWatchID, &input.MediaID)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse query parameters",
			log.ToString("error", err.Error()),
		)
		RespondWithError(log, w, ErrGetMediaWatchInvalidParams, err.Error())
		return
	}
	log.Debug(
		"GetMediaWatch called",
		log.ToString(QueryParamMediaWatchID, strconv.FormatUint(uint64(input.MediaID), 10)),
	)

	// Execute use case
	output, err := h.GetMediaWatchUseCase.Execute(ctx, input)
	if err != nil {
		log.Error(
			"Failed to fetch media",
			log.ToString("error", err.Message),
		)
		RespondWithDTOError(log, w, ErrGetMediaWatchInvalidParams, err)
		return
	}

	log.Debug(
		"GetMediaWatch succeeded",
		log.ToString("media_id", strconv.FormatUint(uint64(input.MediaID), 10)),
		log.ToString("link", output.URL),
	)

	// Respond with output
	Respond(log, w, ResponseGetMediaWatch.Code, output)
}

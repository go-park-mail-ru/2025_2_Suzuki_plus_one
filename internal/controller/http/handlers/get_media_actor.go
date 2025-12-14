package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for GetMediaActor handler
var (
	ErrGetMediaActorInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("invalid parameters for media fetching"),
	}
	ResponseGetMediaActor = Response{
		Code: http.StatusOK,
	}
)

const PathParamGetMediaActorID = "media_id"

// GetMediaActor handler
func (h *Handlers) GetMediaActor(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.GetMediaActorInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddPath(PathParamGetMediaActorID, &input.MediaID)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse query parameters",
			log.ToString("error", err.Error()),
		)
		RespondWithError(log, w, ErrGetMediaActorInvalidParams, err.Error())
		return
	}
	log.Debug(
		"GetMediaActor called for ID",
		log.ToString(PathParamGetMediaActorID, strconv.FormatUint(uint64(input.MediaID), 10)),
	)

	// Execute use case
	output, err := h.GetMediaActorUseCase.Execute(ctx, input)
	if err != nil {
		log.Error(
			"Failed to fetch media",
			log.ToString("error", err.Message),
		)
		RespondWithDTOError(log, w, ErrGetMediaActorInvalidParams, err)
		return
	}

	log.Debug(
		"GetMediaActor succeeded",
		log.ToString("media_id", strconv.FormatUint(uint64(input.MediaID), 10)),
	)

	// Respond with output
	RespondEasyJSON(log, w, ResponseGetMediaActor.Code, output)
}

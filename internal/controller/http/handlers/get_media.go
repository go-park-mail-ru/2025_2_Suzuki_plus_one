package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for GetMedia handler
var (
	ErrGetMediaInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("Invalid parameters for media fetching"),
	}
	ResponseGetMedia = Response{
		Code: http.StatusOK,
	}
)

const PathParamGetMediaID = "media_id"

// GetMedia handler
func (h *Handlers) GetMedia(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContexKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.GetMediaInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddPath(PathParamGetMediaID, &input.MediaID)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse query parameters",
			log.ToString("error", err.Error()),
		)
		RespondWithError(log, w, ErrGetMediaInvalidParams, err.Error())
		return
	}
	log.Debug(
		"GetMedia called for ID",
		log.ToString(PathParamGetMediaID, strconv.FormatUint(uint64(input.MediaID), 10)),
	)

	// Execute use case
	output, err := h.GetMediaUseCase.Execute(ctx, input)
	if err != nil {
		log.Error(
			"Failed to fetch media",
			log.ToString("error", err.Message),
		)
		RespondWithDTOError(log, w, ErrGetMediaInvalidParams, err)
		return
	}

	log.Debug(
		"GetMedia succeeded",
		log.ToString("media_id", strconv.FormatUint(uint64(input.MediaID), 10)),
		log.ToString("actor_count", strconv.Itoa(len(output.Actors))),
	)

	// Respond with output
	Respond(log, w, ResponseGetMedia.Code, output)
}

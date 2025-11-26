package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for GetActor handler
// Blank fields are not used and are filled in the handler
var (
	ErrGetActorInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("Invalid parameters for actor"),
	}
	ErrGetActorNotFound = ResponseError{
		Code: http.StatusBadRequest,
	}
	ResponseGetActor = Response{
		Code: http.StatusOK,
	}
)

// Input path parameter
const PathParamGetActorID = "actor_id"

// Get all movies from database
func (h *Handlers) GetActor(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.GetActorInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddPath(PathParamGetActorID, &input.ActorID)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse query parameters",
			log.ToString("error", err.Error()),
		)
		// Respond with error, if input parameters are invalid
		RespondWithError(log, w, ErrGetActorInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.GetActorUseCase.Execute(ctx, input)
	if err != nil {
		log.Error("Failed to fetch actor",
			log.ToString("error", err.Message),
		)
		// Respond with error, if use case execution fails
		RespondWithDTOError(log, w, ErrGetActorInvalidParams, err)
		return
	}

	// Log successful completion
	log.Debug(
		"Fetching actor details completed successfully",
		log.ToString("actor_id", strconv.FormatUint(uint64(input.ActorID), 10)),
	)

	// Respond with output
	Respond(log, w, ResponseGetActor.Code, output)
}

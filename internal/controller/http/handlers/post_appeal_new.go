package handlers

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for PostAppealNew handler
// Blank fields are not used and are filled in the handler
var (
	ErrPostAppealNewInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("invalid parameters for actor"),
	}
	ResponsePostAppealNew = Response{
		Code: http.StatusOK,
	}
)

// Body parameter names for PostAppealNew handler
var (
	BodyParamPostAppealNewTag     = "tag"
	BodyParamPostAppealNewMessage = "message"
	BodyParamPostAppealNewName    = "name"
)

// Get all movies from database
func (h *Handlers) PostAppealNew(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.PostAppealNewInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddAuthHeader(&input.AccessToken)
	rp.AddBody(BodyParamPostAppealNewTag, &input.Tag)
	rp.AddBody(BodyParamPostAppealNewMessage, &input.Message)
	rp.AddBody(BodyParamPostAppealNewName, &input.Name)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse access token",
			log.ToString("error", err.Error()),
		)
		// Respond with error, if input parameters are invalid
		RespondWithError(log, w, ErrPostAppealNewInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.PostAppealNewUseCase.Execute(ctx, input)
	if err != nil {
		log.Error("Failed to fetch appeals",
			log.ToString("error", err.Message),
		)
		// Respond with error, if use case execution fails
		RespondWithDTOError(log, w, ErrPostAppealNewInvalidParams, err)
		return
	}

	// Log successful completion
	log.Debug(
		"Creating new appeal completed successfully",
		log.ToInt("ID: ", int(output.ID)),
	)

	// Respond with output
	Respond(log, w, ResponsePostAppealNew.Code, output)
}

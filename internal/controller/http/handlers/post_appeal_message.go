package handlers

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for PostAppealMessage handler
// Blank fields are not used and are filled in the handler
var (
	ErrPostAppealMessageInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("Invalid parameters for actor"),
	}
	ResponsePostAppealMessage = Response{
		Code: http.StatusOK,
	}
)

// Body parameter names for PostAppealMessage handler
var (
	BodyParamPostAppealMessageMessage = "message"
)

// Get all movies from database
func (h *Handlers) PostAppealMessage(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.PostAppealMessageInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddAuthHeader(&input.AccessToken)
	// Add path parameter
	rp.AddPath(PathParamGetAppealID, &input.AppealID)
	// Add body parameter
	rp.AddBody(BodyParamPostAppealMessageMessage, &input.Message)


	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse access tokenm",
			log.ToString("error", err.Error()),
		)
		// Respond with error, if input parameters are invalid
		RespondWithError(log, w, ErrPostAppealMessageInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.PostAppealMessageUseCase.Execute(ctx, input)
	if err != nil {
		log.Error("Failed to fetch appeals",
			log.ToString("error", err.Message),
		)
		// Respond with error, if use case execution fails
		RespondWithDTOError(log, w, ErrPostAppealMessageInvalidParams, err)
		return
	}

	// Log successful completion
	log.Debug(
		"Adding a message to specific appeal completed successfully",
		log.ToInt("ID:", int(output.ID)),
	)

	// Respond with output
	Respond(log, w, ResponsePostAppealMessage.Code, output)
}

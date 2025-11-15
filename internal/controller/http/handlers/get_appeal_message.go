package handlers

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for GetAppealMessage handler
// Blank fields are not used and are filled in the handler
var (
	ErrGetAppealMessageInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("Invalid parameters for actor"),
	}
	ResponseGetAppealMessage = Response{
		Code: http.StatusOK,
	}
)

// Get all movies from database
func (h *Handlers) GetAppealMessage(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.GetAppealMessageInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddAuthHeader(&input.AccessToken)
	// Add path parameter
	rp.AddPath(PathParamGetAppealID, &input.AppealId)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse access tokenm",
			log.ToString("error", err.Error()),
		)
		// Respond with error, if input parameters are invalid
		RespondWithError(log, w, ErrGetAppealMessageInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.GetAppealMessageUseCase.Execute(ctx, input)
	if err != nil {
		log.Error("Failed to fetch appeals",
			log.ToString("error", err.Message),
		)
		// Respond with error, if use case execution fails
		RespondWithDTOError(log, w, ErrGetAppealMessageInvalidParams, err)
		return
	}

	// Log successful completion
	log.Debug(
		"Fetching messages of specific appeal completed successfully",
		log.ToInt("messages count", len(output.Messages)),
	)

	// Respond with output
	Respond(log, w, ResponseGetAppealMessage.Code, output)
}

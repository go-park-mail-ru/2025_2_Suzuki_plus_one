package handlers

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for GetAppeal handler
// Blank fields are not used and are filled in the handler
var (
	ErrGetAppealInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("invalid parameters for appeal"),
	}
	ResponseGetAppeal = Response{
		Code: http.StatusOK,
	}
)

// Input path parameter
const PathParamGetAppealID = "appeal_id"

// Get all movies from database
func (h *Handlers) GetAppeal(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.GetAppealInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddAuthHeader(&input.AccessToken)
	rp.AddPath(PathParamGetAppealID, &input.AppealId)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse access tokenm",
			log.ToString("error", err.Error()),
		)
		// Respond with error, if input parameters are invalid
		RespondWithError(log, w, ErrGetAppealInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.GetAppealUseCase.Execute(ctx, input)
	if err != nil {
		log.Error("Failed to fetch appeals",
			log.ToString("error", err.Message),
		)
		// Respond with error, if use case execution fails
		RespondWithDTOError(log, w, ErrGetAppealInvalidParams, err)
		return
	}

	// Log successful completion
	log.Debug(
		"Fetching appeals details completed successfully",
		log.ToInt("ID", int(output.ID)),
	)

	// Respond with output
	RespondEasyJSON(log, w, ResponseGetAppeal.Code, output)
}

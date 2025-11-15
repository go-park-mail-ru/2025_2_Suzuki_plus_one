package handlers

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for PutAppealResolve handler
// Blank fields are not used and are filled in the handler
var (
	ErrPutAppealResolveInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("Invalid parameters for actor"),
	}
	ResponsePutAppealResolve = Response{
		Code: http.StatusOK,
	}
)

// Get all movies from database
func (h *Handlers) PutAppealResolve(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.PutAppealResolveInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddAuthHeader(&input.AccessToken)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse access tokenm",
			log.ToString("error", err.Error()),
		)
		// Respond with error, if input parameters are invalid
		RespondWithError(log, w, ErrPutAppealResolveInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.PutAppealResolveUseCase.Execute(ctx, input)
	if err != nil {
		log.Error("Failed to fetch appeals",
			log.ToString("error", err.Message),
		)
		// Respond with error, if use case execution fails
		RespondWithDTOError(log, w, ErrPutAppealResolveInvalidParams, err)
		return
	}

	// Log successful completion
	log.Debug(
		"Marking an appeal as resolved completed successfully",
		log.ToInt("ID:", int(output.ID)),
	)

	// Respond with output
	Respond(log, w, ResponsePutAppealResolve.Code, output)
}

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
	ErrGetAppealAllInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("Invalid parameters for getting all appeals"),
	}
	ResponseGetAppealAll = Response{
		Code: http.StatusOK,
	}
)

var (
	// URL query parameters
	QueryParamAppealAllOffset = "offset"
	QueryParamAppealAllLimit  = "limit"
	QueryParamAppealAllTag    = "tag"
	QueryParamAppealAllStatus = "status"
)

// Get all movies from database
func (h *Handlers) GetAppealAll(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.GetAppealAllInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddQuery(QueryParamAppealAllTag, &input.Tag)
	rp.AddQuery(QueryParamAppealAllStatus, &input.Status)
	rp.AddQuery(QueryParamAppealAllLimit, &input.Limit)
	rp.AddQuery(QueryParamAppealAllOffset, &input.Offset)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse request parameters",
			log.ToString("error", err.Error()),
		)
		// Respond with error, if input parameters are invalid
		RespondWithError(log, w, ErrGetAppealAllInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.GetAppealAllUseCase.Execute(ctx, input)
	if err != nil {
		log.Error("Failed to fetch appeals",
			log.ToString("error", err.Message),
		)
		// Respond with error, if use case execution fails
		RespondWithDTOError(log, w, ErrGetAppealAllInvalidParams, err)
		return
	}

	// Log successful completion
	log.Debug(
		"Fetching appeals details completed successfully",
		log.ToInt("appeals count", len(output.Appeals)),
	)

	// Respond with output
	Respond(log, w, ResponseGetAppealAll.Code, output)
}

package handlers

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for GetSearch handler
// Blank fields are not used and are filled in the handler
var (
	ErrGetSearchInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("invalid parameters for search"),
	}
	ErrGetSearchNotFound = ResponseError{
		Code: http.StatusBadRequest,
	}
	ResponseGetSearch = Response{
		Code: http.StatusOK,
	}
)

// Input path parameter
var (
	// URL query parameters
	QueryParamSearchQuery  = "query"
	QueryParamSearchOffset = "offset"
	QueryParamSearchLimit  = "limit"
	QueryParamSearchType   = "type"
)

// Get all movies from database
func (h *Handlers) GetSearch(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.GetSearchInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddQuery(QueryParamSearchQuery, &input.Query)
	rp.AddQuery(QueryParamSearchLimit, &input.Limit)
	rp.AddQuery(QueryParamSearchOffset, &input.Offset)
	rp.AddQuery(QueryParamSearchType, &input.Type)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse query parameters",
			log.ToString("error", err.Error()),
		)
		// Respond with error, if input parameters are invalid
		RespondWithError(log, w, ErrGetSearchInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.GetSearchUseCase.Execute(ctx, input)
	if err != nil {
		log.Error("Failed to fetch actor",
			log.ToString("error", err.Message),
		)
		// Respond with error, if use case execution fails
		RespondWithDTOError(log, w, ErrGetSearchInvalidParams, err)
		return
	}

	// Log successful completion
	log.Debug(
		"Fetching search details completed successfully",
		log.ToInt("medias_count", len(output.Medias)),
		log.ToInt("actors_count", len(output.Actors)),
	)

	// Respond with output
	RespondEasyJSON(log, w, ResponseGetSearch.Code, output)
}

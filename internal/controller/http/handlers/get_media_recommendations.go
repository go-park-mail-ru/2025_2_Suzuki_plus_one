package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for GetMediaRecommendations handler
var (
	ErrGetMediaRecommendationsInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("Invalid parameters for movie recommendations"),
	}
	ResponseGetMediaRecommendations = Response{
		Code: http.StatusOK,
	}
)

var (
	// URL query parameters
	QueryParamMovieRecommendationsOffset = "offset"
	QueryParamMovieRecommendationsLimit  = "limit"
	QueryParamMovieRecommendationsType   = "type"
)

// GetMediaRecommendations handler
func (h *Handlers) GetMediaRecommendations(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.GetMediaRecommendationsInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddQuery(QueryParamMovieRecommendationsLimit, &input.Limit)
	rp.AddQuery(QueryParamMovieRecommendationsOffset, &input.Offset)
	rp.AddQuery(QueryParamMovieRecommendationsType, &input.Type)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse query parameters",
			log.ToString("error", err.Error()),
		)
		RespondWithError(log, w, ErrGetMediaRecommendationsInvalidParams, err.Error())
		return
	}
	log.Debug(
		"GetMediaRecommendations called with params",
		log.ToString(QueryParamMovieRecommendationsOffset, strconv.FormatUint(uint64(input.Offset), 10)),
		log.ToString(QueryParamMovieRecommendationsLimit, strconv.FormatUint(uint64(input.Limit), 10)),
		log.ToString(QueryParamMovieRecommendationsType, input.Type),
	)

	// Execute use case
	output, err := h.GetMediaRecommendationsUseCase.Execute(ctx, input)
	if err != nil {
		log.Error(
			"Failed to fetch movie recommendations",
			log.ToString("error", err.Message),
		)
		RespondWithDTOError(log, w, ErrGetMediaRecommendationsInvalidParams, err)
		return
	}

	log.Debug(
		"GetMediaRecommendations succeeded",
		log.ToString("count", strconv.FormatInt(int64(len(output.Movies)), 10)),
		log.ToString("offset", strconv.FormatUint(uint64(input.Offset), 10)),
		log.ToString("limit", strconv.FormatUint(uint64(input.Limit), 10)),
	)

	// Respond with output
	Respond(log, w, ResponseGetMediaRecommendations.Code, output)
}

package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for GetMovieRecommendations handler
var (
	ErrGetMovieRecommendationsInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("Invalid parameters for movie recommendations"),
	}
	ResponseGetMovieRecommendations = Response{
		Code: http.StatusOK,
	}
)

var (
	// URL query parameters
	QueryParamMovieRecommendationsOffset = "offset"
	QueryParamMovieRecommendationsLimit  = "limit"
)

// GetMovieRecommendations handler
func (h *Handlers) GetMovieRecommendations(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContexKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.GetMovieRecommendationsInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddQuery(QueryParamMovieRecommendationsLimit, &input.Limit)
	rp.AddQuery(QueryParamMovieRecommendationsOffset, &input.Offset)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse query parameters",
			log.ToString("error", err.Error()),
		)
		RespondWithError(log, w, ErrGetMovieRecommendationsInvalidParams, err.Error())
		return
	}
	log.Debug(
		"GetMovieRecommendations called with params",
		log.ToString(QueryParamMovieRecommendationsOffset, strconv.FormatUint(uint64(input.Offset), 10)),
		log.ToString(QueryParamMovieRecommendationsLimit, strconv.FormatUint(uint64(input.Limit), 10)),
	)

	// Execute use case
	output, err := h.GetMovieRecommendationsUseCase.Execute(ctx, input)
	if err != nil {
		log.Error(
			"Failed to fetch movie recommendations",
			log.ToString("error", err.Message),
		)
		RespondWithDTOError(log, w, ErrGetMovieRecommendationsInvalidParams, err)
		return
	}

	log.Debug(
		"GetMovieRecommendations succeeded",
		log.ToString("count", strconv.FormatInt(int64(len(output.Movies)), 10)),
		log.ToString("offset", strconv.FormatUint(uint64(input.Offset), 10)),
		log.ToString("limit", strconv.FormatUint(uint64(input.Limit), 10)),
	)

	// Respond with output
	Respond(log, w, ResponseGetMovieRecommendations.Code, output)
}

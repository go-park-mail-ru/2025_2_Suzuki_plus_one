package handlers

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for GetObjectMedia handler
var (
	ErrGetObjectMediaInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("Invalid parameters for media s3 object"),
	}
	ResponseGetObjectMedia = Response{
		Code: http.StatusOK,
	}
)

var (
	// URL query parameters
	QueryParamObjectMediaKey        = "key"
	QueryParamObjectMediaBucketName = "bucket_name"
)

// GetObjectMedia handler
func (h *Handlers) GetObjectMedia(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.GetObjectInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddQuery(QueryParamObjectMediaKey, &input.Key)
	rp.AddQuery(QueryParamObjectMediaBucketName, &input.BucketName)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse query parameters",
			log.ToString("error", err.Error()),
		)
		RespondWithError(log, w, ErrGetObjectMediaInvalidParams, err.Error())
		return
	}
	log.Debug(
		"GetObjectMedia called with params",
		log.ToString(QueryParamObjectMediaKey, input.Key),
		log.ToString(QueryParamObjectMediaBucketName, input.BucketName),
	)

	// Execute use case
	output, err := h.GetObjectMediaUseCase.Execute(ctx, input)
	if err != nil {
		log.Error(
			"Failed to fetch object media",
			log.ToString("error", err.Message),
		)
		RespondWithDTOError(log, w, ErrGetObjectMediaInvalidParams, err)
		return
	}

	log.Debug(
		"GetObjectMedia succeeded",
		log.ToString("url", output.URL),
	)

	// Respond with output
	Respond(log, w, ResponseGetObjectMedia.Code, output)
}

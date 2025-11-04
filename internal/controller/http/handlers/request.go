package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type RequestParams struct {
	logger             logger.Logger
	queryParams        []string
	queryParamsStorage []any
	request            *http.Request
	dto                *dto.DTO
}

func NewRequestParams(l logger.Logger, request *http.Request, dto dto.DTO) *RequestParams {
	return &RequestParams{
		logger:             l,
		queryParams:        make([]string, 0),
		queryParamsStorage: make([]any, 0),
		request:            request,
		dto:                &dto,
	}
}

// Register a query parameter to be parsed with [RequestParams.Parse] into valueStorage (dto field).
// If the parameter is not present in the request, the value in valueStorage becomes zero value.
func (rp *RequestParams) AddQuery(key string, valueStorage any) *RequestParams {
	rp.queryParams = append(rp.queryParams, key)
	rp.queryParamsStorage = append(rp.queryParamsStorage, valueStorage)
	return rp
}

// Parse all registered parameters from the request into their storages.
func (rp *RequestParams) Parse() error {
	// Check request_id exists
	requestID := middleware.GetReqID(rp.request.Context())
	if requestID == "" {
		rp.logger.Error("Can't get requestID",
			rp.logger.ToString("requestURI", rp.request.URL.String()))
		return errors.New("no request ID found: " + rp.request.URL.String())
	}

	query := rp.request.URL.Query()
	for i := range rp.queryParams {
		param := rp.queryParams[i]
		storage := rp.queryParamsStorage[i]
		val := query.Get(param)

		// Scan value into the given storage
		if _, err := fmt.Sscanf(val, "%v", storage); err != nil {
			rp.logger.Warn("Invalid query parameter",
				rp.logger.ToString("param", param),
				rp.logger.ToString("value", val),
				rp.logger.ToError(err))
			// If can't scan just set zero value
		}
	}
	return nil
}

func (rp *RequestParams) GetContext() context.Context {
	// request_id is injected by chi, but we don't want to depend on chi everywhere
	// so we create our own context with common.RequestIDContextKey
	ctx := context.WithValue(rp.request.Context(), common.RequestIDContextKey, middleware.GetReqID(rp.request.Context()))
	return ctx
}

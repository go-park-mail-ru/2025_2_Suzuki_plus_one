package handlers

import (
	"context"
	"encoding/json"
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
	bodyParams         []string
	bodyParamsStorage  []any
	request            *http.Request
	dto                *dto.DTO
}

func NewRequestParams(l logger.Logger, request *http.Request, dto dto.DTO) *RequestParams {
	return &RequestParams{
		logger:             l,
		queryParams:        make([]string, 0),
		queryParamsStorage: make([]any, 0),
		bodyParams:         make([]string, 0),
		bodyParamsStorage:  make([]any, 0),
		request:            request,
		dto:                &dto,
	}
}

// Register a query parameter to be parsed with [RequestParams.Parse] into valueStorage (dto field).
// If the parameter is not present in the request, the value in valueStorage becomes zero value.
func (rp *RequestParams) AddQuery(key string, valueStorage any) {
	rp.queryParams = append(rp.queryParams, key)
	rp.queryParamsStorage = append(rp.queryParamsStorage, valueStorage)
}

func (rp *RequestParams) AddBody(key string, valueStorage any) {
	rp.bodyParams = append(rp.bodyParams, key)
	rp.bodyParamsStorage = append(rp.bodyParamsStorage, valueStorage)
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

	// Read and parse body parameters if any
	if len(rp.bodyParams) == 0 {
		return nil
	}

	// ioBody, err := io.ReadAll(rp.request.Body)
	// if err != nil {
	// 	rp.logger.Error("Failed to read request body",
	// 		rp.logger.ToString("requestURI", rp.request.URL.String()),
	// 		rp.logger.ToError(err))
	// 	return err
	// }
	// fmt.Println(string(ioBody))
	// if err := json.Unmarshal(ioBody, &rp.bodyParamsStorage); err != nil {
	// 	rp.logger.Warn("Failed to decode body parameters",
	// 		rp.logger.ToError(err))
	// 	return err
	// }
	if err := json.NewDecoder(rp.request.Body).Decode(rp.dto); err != nil {
		rp.logger.Warn("Failed to decode body parameters",
			rp.logger.ToError(err))
		return err
	}

	return nil
}

func (rp *RequestParams) GetContext() context.Context {
	// request_id is injected by chi, but we don't want to depend on chi everywhere
	// so we create our own context with common.RequestIDContextKey
	ctx := context.WithValue(rp.request.Context(), common.RequestIDContextKey, middleware.GetReqID(rp.request.Context()))
	return ctx
}

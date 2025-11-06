package handlers

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for GetAuthRefresh handler
var (
	ErrGetAuthRefreshInvalidParams = ResponseError{
		Code:    http.StatusUnauthorized,
		Message: errors.New("Invalid parameters for auth refresh"),
	}
	ResponseGetAuthRefresh = Response{
		Code: http.StatusOK,
	}
)

// Cookie input parameter
var CookieParamGetAuthRefresh = CookieNameRefreshToken

// GetAuthRefresh handler
func (h *Handlers) GetAuthRefresh(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContexKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.GetAuthRefreshInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddCookie(CookieParamGetAuthRefresh, &input.RefreshToken)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse query parameters",
			log.ToString("error", err.Error()),
		)
		RespondWithError(log, w, ErrGetAuthRefreshInvalidParams, err.Error())
		return
	}
	log.Debug(
		"GetAuthRefresh called",
		log.ToString(CookieParamGetAuthRefresh, input.RefreshToken),
	)

	// Execute use case
	output, err := h.GetAuthRefreshUseCase.Execute(ctx, input)
	if err != nil {
		log.Error(
			"Failed to refresh auth",
			log.ToString("error", err.Message),
		)
		RespondWithDTOError(log, w, ErrGetAuthRefreshInvalidParams, err)
		return
	}

	log.Debug(
		"GetAuthRefresh succeeded",
		log.ToString("refresh_token", input.RefreshToken),
		log.ToString("access_token", output.AccessToken),
	)

	// Respond with output
	Respond(log, w, ResponseGetAuthRefresh.Code, output)
}

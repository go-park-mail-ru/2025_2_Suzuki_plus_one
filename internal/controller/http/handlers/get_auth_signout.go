package handlers

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for GetAuthSignOut handler
var (
	ErrGetAuthSignOutInvalidParams = ResponseError{
		Code:    http.StatusUnauthorized,
		Message: errors.New("Invalid parameters for authentication sign-out"),
	}
	ResponseGetAuthSignOut = Response{
		Code: http.StatusOK,
	}
)

// Cookie input parameter
var CookieParamGetAuthSignOut = CookieNameRefreshToken

// Cookie that will be set in output
var CookieOutputGetAuthSignOut = CookieNameRefreshToken

// GetAuthSignOut handler
func (h *Handlers) GetAuthSignOut(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContexKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.GetAuthSignOutInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddCookie(CookieParamGetAuthSignOut, &input.RefreshToken)
	rp.AddAuthHeader(&input.AccessToken)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse token parameters",
			log.ToString("error", err.Error()),
		)
		RespondWithError(log, w, ErrGetAuthSignOutInvalidParams, err.Error())
		return
	}
	log.Debug(
		"GetAuthSignOut called",
		log.ToString(CookieParamGetAuthSignOut, input.RefreshToken),
		log.ToString("access_token", input.AccessToken),
	)

	// Execute use case
	output, err := h.GetAuthSignOutUseCase.Execute(ctx, input)
	if err != nil {
		log.Error(
			"Failed to sign out user",
			log.ToString("error", err.Message),
		)
		RespondWithDTOError(log, w, ErrGetAuthSignOutInvalidParams, err)
		return
	}

	// Reset refresh token in cookie
	resetCookie := NewResetCookieRefreshToken()
	http.SetCookie(w, resetCookie)

	log.Debug(
		"GetAuthSignOut succeeded",
		log.ToString("refresh_token", input.RefreshToken),
		log.ToString("access_token", input.AccessToken),
	)

	// Respond with output
	Respond(log, w, ResponseGetAuthSignOut.Code, output)
}

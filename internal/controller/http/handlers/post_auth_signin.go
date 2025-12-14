package handlers

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for PostAuthSignIn handler
var (
	ErrPostAuthSignInInvalidParams = ResponseError{
		Code:    http.StatusUnauthorized,
		Message: errors.New("invalid parameters for authentication sign-in"),
	}
	ResponsePostAuthSignIn = Response{
		Code: http.StatusOK,
	}
)

// Cookie input parameter
var CookieParamPostAuthSignIn = CookieRefreshTokenName

// Cookie that will be set in output
var CookieOutputPostAuthSignIn = CookieRefreshTokenName

// PostAuthSignIn handler
func (h *Handlers) PostAuthSignIn(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.PostAuthSignInInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddBody("email", &input.Email)
	rp.AddBody("password", &input.Password)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse body parameters",
			log.ToString("error", err.Error()),
		)
		RespondWithError(log, w, ErrPostAuthSignInInvalidParams, err.Error())
		return
	}
	log.Debug(
		"PostAuthSignIn called",
		log.ToString("email", input.Email),
	)

	// Execute use case
	output, err := h.PostAuthSignInUseCase.Execute(ctx, input)
	if err != nil {
		log.Error(
			"Failed to sign in user",
			log.ToString("error", err.Message),
		)
		RespondWithDTOError(log, w, ErrPostAuthSignInInvalidParams, err)
		return
	}

	// Set refresh token in cookie
	refreshTokenCookie := NewRefreshTokenCookie(output.RefreshToken)
	http.SetCookie(w, refreshTokenCookie)

	log.Debug(
		"PostAuthSignIn succeeded",
		log.ToString("email", input.Email),
		log.ToString("access_token", output.AccessToken),
		log.ToString("refresh_token", output.RefreshToken),
	)

	// Respond with output (refresh token not in body)
	RespondEasyJSON(log, w, ResponsePostAuthSignIn.Code, UpdatePostAuthSignInOutput(output))
}

// UpdatePostAuthSignInOutput modifies the output to avoid sending the refresh token in the response body
func UpdatePostAuthSignInOutput(input dto.PostAuthSignInOutput) dto.PostAuthSignInOutput {
	return dto.PostAuthSignInOutput{
		AccessToken:  input.AccessToken,
		RefreshToken: "token is set in http-only cookie",
	}
}

package handlers

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for PostAuthSignUp handler
var (
	ErrPostAuthSignUpInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("Invalid parameters for post-auth signup"),
	}
	ResponsePostAuthSignUp = Response{
		Code: http.StatusOK,
	}
)

// Cookie input parameter
var CookieParamPostAuthSignUp = CookieNameRefreshToken

// Cookie that will be set in output
var CookieOutputPostAuthSignUp = CookieNameRefreshToken

// PostAuthSignUp handler
func (h *Handlers) PostAuthSignUp(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContexKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.PostAuthSignUpInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddBody("email", &input.Email)
	rp.AddBody("password", &input.Password)
	rp.AddBody("username", &input.Username)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse signup parameters",
			log.ToString("error", err.Error()),
		)
		RespondWithError(log, w, ErrPostAuthSignUpInvalidParams, err.Error())
		return
	}
	log.Debug(
		"PostAuthSignUp called",
		log.ToString("email", input.Email),
		log.ToString("username", input.Username),
	)

	// Execute use case
	output, err := h.PostAuthSignUpUseCase.Execute(ctx, input)
	if err != nil {
		log.Error(
			"Failed to sign up user",
			log.ToString("error", err.Message),
		)
		RespondWithDTOError(log, w, ErrPostAuthSignUpInvalidParams, err)
		return
	}

	// Set refresh token in cookie
	refreshTokenCookie := NewRefreshTokenCookie(output.RefreshToken)
	http.SetCookie(w, refreshTokenCookie)

	log.Debug(
		"PostAuthSignUp succeeded",
		log.ToString("access_token", output.AccessToken),
		log.ToString("refresh_token", output.RefreshToken),
	)

	// Respond with output
	Respond(log, w, ResponsePostAuthSignUp.Code, UpdatePostAuthSignUpOutput(output))
}

// UpdatePostAuthSignUpOutput modifies the output to avoid sending the refresh token in the response body
func UpdatePostAuthSignUpOutput(input dto.PostAuthSignUpOutput) dto.PostAuthSignUpOutput {
	return dto.PostAuthSignUpOutput{
		AccessToken:  input.AccessToken,
		RefreshToken: "token is set in http-only cookie",
	}
}

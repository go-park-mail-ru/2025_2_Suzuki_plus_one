package handlers

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for PostUserMeUpdatePassword handler
// Blank fields are not used and are filled in the handler
var (
	ErrPostUserMeUpdatePasswordInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("invalid parameters for user password update"),
	}
	ErrPostUserMeUpdatePasswordNotFound = ResponseError{
		Code: http.StatusBadRequest,
	}
	ResponsePostUserMeUpdatePassword = Response{
		Code: http.StatusOK,
	}
)

// Body parameter names for PostUserMeUpdatePassword handler
var (
	BodyParamPostUserMeUpdatePasswordCurrentPassword = "current_password"
	BodyParamPostUserMeUpdatePasswordNewPassword     = "new_password"
)

// Get all movies from database
func (h *Handlers) PostUserMeUpdatePassword(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.PostUserMeUpdatePasswordInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddAuthHeader(&input.AccessToken)

	// Bind body parameters.
	// Note: this is unnecessary,
	// because we parse body with json tags in [dto.PostUserMeUpdatePasswordInput]
	rp.AddBody(BodyParamPostUserMeUpdatePasswordCurrentPassword, &input.CurrentPassword)
	rp.AddBody(BodyParamPostUserMeUpdatePasswordNewPassword, &input.NewPassword)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse parameters",
			log.ToString("error", err.Error()),
		)
		// Respond with error, if input parameters are invalid
		RespondWithError(log, w, ErrPostUserMeUpdatePasswordInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.PostUserMeUpdatePasswordUseCase.Execute(ctx, input)
	if err != nil {
		log.Error("Failed to fetch actor",
			log.ToString("error", err.Message),
		)
		// Respond with error, if use case execution fails
		RespondWithDTOError(log, w, ErrPostUserMeUpdatePasswordInvalidParams, err)
		return
	}

	// Log successful completion
	log.Debug(
		"Password update successful",
	)

	// Respond with output
	Respond(log, w, ResponsePostUserMeUpdatePassword.Code, output)
}

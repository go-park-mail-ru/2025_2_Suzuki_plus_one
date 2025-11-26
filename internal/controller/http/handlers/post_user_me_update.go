package handlers

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for PostUserMeUpdate handler
// Blank fields are not used and are filled in the handler
var (
	ErrPostUserMeUpdateInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("Invalid parameters for user update"),
	}
	ErrPostUserMeUpdateCantUpdate = ResponseError{
		Code: http.StatusBadRequest,
	}
	ResponsePostUserMeUpdate = Response{
		Code: http.StatusOK,
	}
)

// Get all movies from database
func (h *Handlers) PostUserMeUpdate(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.PostUserMeUpdateInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddAuthHeader(&input.AccessToken)
	rp.AddBody("email", &input.Email)
	rp.AddBody("username", &input.Username)
	rp.AddBody("date_of_birth", &input.DateOfBirth)
	rp.AddBody("phone_number", &input.PhoneNumber)

	// Parse request BODY parameters into the input DTO
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse body parameters",
			log.ToString("error", err.Error()),
		)
		// Respond with error, if input parameters are invalid
		RespondWithError(log, w, ErrPostUserMeUpdateInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.PostUserMeUpdateUseCase.Execute(ctx, input)
	if err != nil {
		log.Error("Failed to update user profile",
			log.ToString("error", err.Message),
		)
		// Respond with error, if use case execution fails
		RespondWithDTOError(log, w, ErrPostUserMeUpdateInvalidParams, err)
		return
	}

	// Log successful completion
	log.Debug(
		"Updating user profile completed successfully",
		log.ToString("email", input.Email),
	)

	// Respond with output
	Respond(log, w, ResponsePostUserMeUpdate.Code, output)
}

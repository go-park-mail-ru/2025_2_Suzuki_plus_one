package handlers

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for GetUserMe handler
var (
	ErrGetUserMeInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("Invalid parameters for getting user info"),
	}
	ResponseGetUserMe = Response{
		Code: http.StatusOK,
	}
)

// GetUserMe handler
func (h *Handlers) GetUserMe(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContexKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.GetUserMeInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddAuthHeader(&input.AccessToken)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse token parameters",
			log.ToString("error", err.Error()),
		)
		RespondWithError(log, w, ErrGetUserMeInvalidParams, err.Error())
		return
	}
	log.Debug(
		"GetUserMe called with params",
		log.ToString("access_token", input.AccessToken),
	)

	// Execute use case
	output, err := h.GetUserMeUseCase.Execute(ctx, input)
	if err != nil {
		log.Error(
			"Failed to get user info",
			log.ToString("error", err.Message),
		)
		RespondWithDTOError(log, w, ErrGetUserMeInvalidParams, err)
		return
	}

	log.Debug(
		"GetUserMe succeeded for user ID",
		log.ToInt("user_id", int(output.ID)),
	)

	// Respond with output
	Respond(log, w, ResponseGetUserMe.Code, output)
}

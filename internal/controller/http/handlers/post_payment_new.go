package handlers

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for PostPaymentNew handler
// Blank fields are not used and are filled in the handler
var (
	ErrPostPaymentNewInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("Invalid parameters for payment"),
	}
	ResponsePostPaymentNew = Response{
		Code: http.StatusCreated,
	}
)

// Get all movies from database
func (h *Handlers) PostPaymentNew(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.PostPaymentNewInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddAuthHeader(&input.AccessToken)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse access token",
			log.ToString("error", err.Error()),
		)
		// Respond with error, if input parameters are invalid
		RespondWithError(log, w, ErrPostPaymentNewInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.PostPaymentNewUseCase.Execute(ctx, input)
	if err != nil {
		log.Error("Failed to create new payment",
			log.ToString("error", err.Message),
		)
		// Respond with error, if use case execution fails
		RespondWithDTOError(log, w, ErrPostPaymentNewInvalidParams, err)
		return
	}

	// Log successful completion
	log.Debug(
		"Creating new payment completed successfully",
	)

	// Respond with redirect to confirmation URL
	http.Redirect(w, r, output.ConfirmationURL, http.StatusSeeOther)
}

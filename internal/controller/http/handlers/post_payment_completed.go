package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for PostPaymentCompleted handler
// Blank fields are not used and are filled in the handler
var (
	ErrPostPaymentCompletedInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("invalid parameters for actor"),
	}
	ErrPostPaymentCompletedNotFound = ResponseError{
		Code: http.StatusBadRequest,
	}
	ResponsePostPaymentCompleted = Response{
		Code: http.StatusOK,
	}
)

// Get all movies from database
func (h *Handlers) PostPaymentCompleted(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.PostPaymentCompletedInput{}
	rp := NewRequestParams(log, r, &input)

	// TODO: use rp.addbody?
	// Decode webhook event from request body
	if err := json.NewDecoder(r.Body).Decode(&input.Webhook); err != nil {
		log.Error(
			"Failed to decode webhook event",
			log.ToError(err),
		)
		// Respond with error, if input parameters are invalid
		RespondWithError(log, w, ErrPostPaymentCompletedInvalidParams, err.Error())
		return
	}
	log.Info("Received webhook event",
		log.ToAny("webhookEvent", input.Webhook),
	)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse query parameters",
			log.ToString("error", err.Error()),
		)
		// Respond with error, if input parameters are invalid
		RespondWithError(log, w, ErrPostPaymentCompletedInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.PostPaymentCompletedUseCase.Execute(ctx, input)
	if err != nil {
		log.Error("Failed to fetch actor",
			log.ToString("error", err.Message),
		)
		// Respond with error, if use case execution fails
		RespondWithDTOError(log, w, ErrPostPaymentCompletedInvalidParams, err)
		return
	}

	// Log successful completion
	log.Debug(
		"Fetching actor details completed successfully",
	)

	// Respond with output
	Respond(log, w, ResponsePostPaymentCompleted.Code, output)
}

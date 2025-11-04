package handlers

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
)

var (
	ErrObjectMediaInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: "Invalid parameters for media s3 object",
	}
)

// Get all media objects from database
func (h *Handlers) GetObjectMedia(w http.ResponseWriter, r *http.Request) {
	// Handle input parameters
	input := dto.GetObjectInput{}
	rp := NewRequestParams(h.logger, r, &input)
	rp.AddQuery("key", &input.Key)
	rp.AddQuery("bucket_name", &input.BucketName)
	if err := rp.Parse(); err != nil {
		h.logger.Error("Failed to parse query parameters",
			h.logger.ToString("error", err.Error()))
		h.ResponseWithError(w, ErrObjectMediaInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.GetObjectMediaUseCase.Execute(rp.GetContext(), input)
	if err != nil {
		h.logger.Error("Failed to fetch movie recommendations",
			h.logger.ToString("error", err.Message))
		// Respond with error
		h.Response(w, ErrObjectMediaInvalidParams.Code, err)
		return
	}

	h.logger.Info("Fetching movie recommendations completed successfully",
		h.logger.ToString("url", output.URL))

	// Respond with output
	h.Response(w, http.StatusOK, output)
}

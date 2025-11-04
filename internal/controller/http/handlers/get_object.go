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
	rp := NewRequestParams(h.Logger, r, &input)
	rp.AddQuery("key", &input.Key)
	rp.AddQuery("bucket_name", &input.BucketName)
	if err := rp.Parse(); err != nil {
		h.Logger.Error("Failed to parse query parameters",
			h.Logger.ToString("error", err.Error()))
		h.ResponseWithError(w, ErrObjectMediaInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.GetObjectMediaUseCase.Execute(rp.GetContext(), input)
	if err != nil {
		h.Logger.Error("Failed to fetch movie recommendations",
			h.Logger.ToString("error", err.Message))
		// Respond with error
		h.Response(w, ErrObjectMediaInvalidParams.Code, err)
		return
	}

	h.Logger.Info("Fetching movie recommendations completed successfully",
		h.Logger.ToString("url", output.URL))

	// Respond with output
	h.Response(w, http.StatusOK, output)
}

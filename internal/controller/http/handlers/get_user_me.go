package handlers

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
)

var (
	ErrGetUserMeInvalidParams = ResponseError{
		http.StatusBadRequest,
		"Invalid parameters for getting user info",
	}
)

func (h *Handlers) GetUserMe(w http.ResponseWriter, r *http.Request) {
	input := dto.GetUserMeInput{}

	// Parse Access Token from header
	rp := NewRequestParams(h.Logger, r, &input)
	rp.AddAuthHeader(&input.AccessToken)
	if err := rp.Parse(); err != nil {
		h.Logger.Error("Failed to parse token parameters",
			h.Logger.ToString("error", err.Error()))
		h.ResponseWithError(w, ErrGetUserMeInvalidParams, err.Error())
		return
	}
	h.Logger.Debug("GetUserMe called",
		h.Logger.ToString("access_token", input.AccessToken),
	)

	// Execute use case
	output, err := h.GetUserMeUseCase.Execute(r.Context(), input)
	if err != nil {
		h.Logger.Error("Failed to get user info",
			h.Logger.ToString("error", err.Message))
		h.Response(w, ErrGetUserMeInvalidParams.Code, err)
		return
	}

	h.Response(w, http.StatusOK, output)
}

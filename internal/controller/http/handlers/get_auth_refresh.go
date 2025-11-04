package handlers

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
)

var (
	ErrGetAuthRefreshInvalidParams = ResponseError{
		Code:    http.StatusUnauthorized,
		Message: "Invalid parameters for auth refresh",
	}
)

var RefreshTokenCookieName = "refresh_token"

// Get all movies from database
func (h *Handlers) GetAuthRefresh(w http.ResponseWriter, r *http.Request) {
	// Handle input parameters
	input := dto.GetAuthRefreshInput{}
	rp := NewRequestParams(h.Logger, r, &input)
	rp.AddCookie(RefreshTokenCookieName, &input.RefreshToken)
	if err := rp.Parse(); err != nil {
		h.Logger.Error("Failed to parse query parameters",
			h.Logger.ToString("error", err.Error()))
		h.ResponseWithError(w, ErrGetAuthRefreshInvalidParams, err.Error())
		return
	}
	h.Logger.Info("GetAuthRefresh called",
		h.Logger.ToString(RefreshTokenCookieName, input.RefreshToken),
	)

	// Execute use case
	output, err := h.GetAuthRefreshUseCase.Execute(rp.GetContext(), input)
	if err != nil {
		h.Logger.Error("Failed to fetch movie recommendations",
			h.Logger.ToString("error", err.Message))
		// Respond with error
		h.Response(w, ErrGetAuthRefreshInvalidParams.Code, err)
		return
	}

	h.Logger.Info("GetAuthRefresh succeeded",
		h.Logger.ToString("refresh_token", input.RefreshToken),
		h.Logger.ToString("access_token", output.AccessToken),
	)

	// Respond with output
	h.Response(w, http.StatusOK, output)
}

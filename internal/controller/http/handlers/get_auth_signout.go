package handlers

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
)

var (
	ErrGetAuthSignOutInvalidParams = ResponseError{
		Code:    http.StatusUnauthorized,
		Message: "Invalid parameters for authentication sign-out",
	}
)

func NewResetRefreshTokenCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // TODO: set to true in production
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
	}
}

// Sign in user, respond with auth refresh and access tokens
func (h *Handlers) GetAuthSignOut(w http.ResponseWriter, r *http.Request) {
	input := dto.GetAuthSignOutInput{}

	// Parse Refresh Token from cookie and Access Token from header
	rp := NewRequestParams(h.Logger, r, &input)
	rp.AddCookie(RefreshTokenCookieName, &input.RefreshToken)
	rp.AddAuthHeader(&input.AccessToken)
	if err := rp.Parse(); err != nil {
		h.Logger.Error("Failed to parse token parameters",
			h.Logger.ToString("error", err.Error()))
		h.ResponseWithError(w, ErrGetAuthSignOutInvalidParams, err.Error())
		return
	}
	h.Logger.Debug("GetAuthSignOut called",
		h.Logger.ToString(RefreshTokenCookieName, input.RefreshToken),
		h.Logger.ToString("access_token", input.AccessToken),
	)

	// Execute use case
	output, err := h.GetAuthSignOutUseCase.Execute(r.Context(), input)
	if err != nil {
		h.Logger.Error("Failed to sign out user",
			h.Logger.ToString("error", err.Message))
		h.Response(w, ErrGetAuthSignOutInvalidParams.Code, err)
		return
	}

	// Reset refresh token in cookie
	ResetRefreshTokenCookie := NewResetRefreshTokenCookie()
	http.SetCookie(w, ResetRefreshTokenCookie)

	// Respond with output
	h.Response(w, http.StatusOK, output)
}

package handlers

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
)

var (
	ErrAuthSignInInvalidParams = ResponseError{
		Code:    http.StatusUnauthorized,
		Message: "Invalid parameters for authentication sign-in",
	}
)

func NewRefreshTokenCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     "refresh_token",
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,                              // TODO: set to true in production
		Expires:  time.Now().Add(7 * 24 * time.Hour), // TODO: change from 7 days if needed
		SameSite: http.SameSiteStrictMode,
	}
}

// UpdatePostAuthSignInOutput modifies the output to avoid sending the refresh token in the response body
func UpdatePostAuthSignInOutput(input dto.PostAuthSignInOutput) dto.PostAuthSignInOutput {
	return dto.PostAuthSignInOutput{
		AccessToken:  input.AccessToken,
		RefreshToken: "token is set in http-only cookie",
	}
}

// Sign in user, respond with auth refresh and access tokens
func (h *Handlers) PostAuthSignIn(w http.ResponseWriter, r *http.Request) {
	input := dto.PostAuthSignInInput{}
	rp := NewRequestParams(h.Logger, r, &input)
	rp.AddBody("email", &input.Email)
	rp.AddBody("password", &input.Password)
	if err := rp.Parse(); err != nil {
		h.Logger.Error("Failed to parse body parameters",
			h.Logger.ToString("error", err.Error()))
		h.ResponseWithError(w, ErrAuthSignInInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.PostAuthSignInUseCase.Execute(rp.GetContext(), input)
	if err != nil {
		h.Logger.Error("Failed to sign in user",
			h.Logger.ToString("error", err.Message))
		h.Response(w, ErrAuthSignInInvalidParams.Code, err)
		return
	}

	h.Logger.Info("User signed in successfully",
		h.Logger.ToString("email", input.Email),
		h.Logger.ToString("accessToken", output.AccessToken),
		h.Logger.ToString("refreshToken", output.RefreshToken),
	)

	// Set refresh token in cookie
	RefreshTokenCookie := NewRefreshTokenCookie(output.RefreshToken)
	http.SetCookie(w, RefreshTokenCookie)

	// Respond with output
	h.Response(w, http.StatusOK, UpdatePostAuthSignInOutput(output))
}

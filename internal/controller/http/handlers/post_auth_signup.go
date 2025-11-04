package handlers

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
)

var (
	ErrPostAuthSignUpInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: "Invalid parameters for post-auth signup",
	}
)

// UpdatePostAuthSignUpOutput modifies the output to avoid sending the refresh token in the response body
func UpdatePostAuthSignUpOutput(input dto.PostAuthSignUpOutput) dto.PostAuthSignUpOutput {
	return dto.PostAuthSignUpOutput{
		AccessToken:  input.AccessToken,
		RefreshToken: "token is set in http-only cookie",
	}
}

// Get all media objects from database
func (h *Handlers) PostAuthSignUp(w http.ResponseWriter, r *http.Request) {
	// Handle input parameters
	input := dto.PostAuthSignUpInput{}
	rp := NewRequestParams(h.Logger, r, &input)
	rp.AddBody("email", &input.Email)
	rp.AddBody("password", &input.Password)
	rp.AddBody("username", &input.Username)
	if err := rp.Parse(); err != nil {
		h.Logger.Error("Failed to parse query parameters",
			h.Logger.ToString("error", err.Error()))
		h.ResponseWithError(w, ErrPostAuthSignUpInvalidParams, err.Error())
		return
	}

	// Execute use case
	output, err := h.PostAuthSignUpUseCase.Execute(rp.GetContext(), input)
	if err != nil {
		h.Logger.Error("Failed to fetch from database",
			h.Logger.ToString("error", err.Message))
		// Respond with error
		h.Response(w, ErrPostAuthSignUpInvalidParams.Code, err)
		return
	}

	h.Logger.Debug("Signed up user successfully",
		h.Logger.ToString("accessToken", output.AccessToken),
		h.Logger.ToString("refreshToken", output.RefreshToken),
	)

	// Set refresh token in cookie
	RefreshTokenCookie := NewRefreshTokenCookie(output.RefreshToken)
	http.SetCookie(w, RefreshTokenCookie)

	// Respond with output
	h.Response(w, http.StatusOK, UpdatePostAuthSignUpOutput(output))
}

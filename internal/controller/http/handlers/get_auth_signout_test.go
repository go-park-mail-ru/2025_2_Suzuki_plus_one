package handlers_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	srv "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http"
	. "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/handlers"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	. "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/mocks"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
)

func getMockGetAuthSignOutInput() dto.GetAuthSignOutInput {
	token, err := common.GenerateToken(1, time.Hour)
	if err != nil {
		panic("failed to generate access token")
	}
	return dto.GetAuthSignOutInput{
		RefreshToken: "RefreshTokenValue",
		AccessToken:  token,
	}
}

func getMockGetAuthSignOutOutput() dto.GetAuthSignOutOutput {
	return dto.GetAuthSignOutOutput{}
}

// Call GetAuthSignOut handler and check response with query parameters limit and offset
func TestGetAuthSignOut(t *testing.T) {
	logger := logger.NewZapLogger(true)
	common.InitJWT("secret", time.Hour, time.Hour)

	// Define input and expected output
	input := getMockGetAuthSignOutInput()
	output := getMockGetAuthSignOutOutput()

	// Create mock GetAuthSignOutUsecase
	mockCtrl := gomock.NewController(t)
	mockGetAuthSignOutUsecase := NewMockGetAuthSignOutUseCase(mockCtrl)
	mockGetAuthSignOutUsecase.EXPECT().
		Execute(gomock.Any(), gomock.Eq(input)).
		Return(output, nil).
		Times(1)

	// Initialize server with the mock usecase
	handlers := &Handlers{
		Logger:                logger,
		GetAuthSignOutUseCase: mockGetAuthSignOutUsecase,
	}
	router := srv.InitRouter(handlers, logger, "/")
	server := srv.NewServer(router)

	// Create a New Request
	requestURL := "/auth/signout"
	req, err := http.NewRequest("GET", requestURL, nil)
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  CookieNameRefreshToken,
		Value: input.RefreshToken,
	})
	req.Header.Add("Authorization", "Bearer "+input.AccessToken)

	// Execute Request
	response := executeRequest(logger, server, req)

	// Assert the response
	checkResponse(t, logger, response, http.StatusOK, output)
}

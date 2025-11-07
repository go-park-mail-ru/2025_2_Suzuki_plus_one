package handlers_test

import (
	"net/http"
	"testing"

	srv "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http"
	. "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/handlers"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/stretchr/testify/require"

	. "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/mocks"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"

	gomock "go.uber.org/mock/gomock"
)

func getMockGetAuthRefreshInput() dto.GetAuthRefreshInput {
	return dto.GetAuthRefreshInput{
		RefreshToken: "RefreshTokenValue",
	}
}

func getMockGetAuthRefreshOutput() dto.GetAuthRefreshOutput {
	return dto.GetAuthRefreshOutput{
		AccessToken: "AccessTokenValue",
	}
}

// Call GetAuthRefresh handler and check response with query parameters limit and offset
func TestGetAuthRefresh(t *testing.T) {
	logger := logger.NewZapLogger(true)

	// Define input and expected output
	input := getMockGetAuthRefreshInput()
	movies := getMockGetAuthRefreshOutput()

	// Create mock GetAuthRefreshUsecase
	mockCtrl := gomock.NewController(t)
	mockGetAuthRefreshUsecase := NewMockGetAuthRefreshUseCase(mockCtrl)
	mockGetAuthRefreshUsecase.EXPECT().
		Execute(gomock.Any(), gomock.Eq(input)).
		Return(movies, nil).
		Times(1)

	// Initialize server with the mock usecase
	handlers := &Handlers{
		Logger:                logger,
		GetAuthRefreshUseCase: mockGetAuthRefreshUsecase,
	}
	router := srv.InitRouter(handlers, logger, "/")
	server := srv.NewServer(router)

	// Create a New Request
	requestURL := "/auth/refresh"
	req, err := http.NewRequest("GET", requestURL, nil)
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  CookieRefreshTokenName,
		Value: input.RefreshToken,
	})

	// Execute Request
	response := executeRequest(logger, server, req)

	// Assert the response
	checkResponse(t, logger, response, http.StatusOK, movies)
}

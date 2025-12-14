package handlers_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	srv "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http"
	. "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/handlers"
	. "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/mocks"
	rtr "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/router"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/stretchr/testify/require"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	gomock "go.uber.org/mock/gomock"
)

func getMockGetUserMeInput() dto.GetUserMeInput {
	x, err := common.GenerateToken(1, time.Minute)
	if err != nil {
		panic("failed to generate token")
	}
	return dto.GetUserMeInput{
		AccessToken: x,
	}
}

func getMockGetUserMeOutput() dto.GetUserMeOutput {
	return dto.GetUserMeOutput{}
}

// Call GetUserMe handler and check response with query parameters limit and offset
func TestGetUserMe(t *testing.T) {
	logger := logger.NewZapLogger(true)
	common.InitJWT("test_secret", time.Minute, time.Hour)

	// Define input and expected output
	input := getMockGetUserMeInput()
	output := getMockGetUserMeOutput()

	// Create mock GetUserMeUsecase
	mockCtrl := gomock.NewController(t)
	mockGetUserMeUsecase := NewMockGetUserMeUseCase(mockCtrl)
	mockGetUserMeUsecase.EXPECT().
		Execute(gomock.Any(), gomock.Eq(input)).
		Return(output, nil).
		Times(1)

	// Initialize server with the mock usecase
	handlers := &Handlers{
		Logger:           logger,
		GetUserMeUseCase: mockGetUserMeUsecase,
	}
	router := rtr.InitRouter(handlers, logger, "/")
	server := srv.NewServer(router)

	// Create a New Request
	requestURL := "/user/me"
	req, err := http.NewRequest("GET", requestURL, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", input.AccessToken))

	// Execute Request
	response := executeRequest(logger, server, req)

	// Assert the response
	checkResponse(t, logger, response, http.StatusOK, output)
}

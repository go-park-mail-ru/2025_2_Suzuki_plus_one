package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	srv "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http"
	. "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/handlers"
	rtr "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/router"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/stretchr/testify/require"

	. "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/mocks"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	gomock "go.uber.org/mock/gomock"
)

func getMockGetMediaInput() dto.GetMediaInput {
	return dto.GetMediaInput{
		MediaID: 2,
	}
}

func getMockGetMediaOutput() dto.GetMediaOutput {
	return dto.GetMediaOutput{}
}

// Call GetMedia handler and check response with query parameters limit and offset
func TestGetMedia(t *testing.T) {
	logger := logger.NewZapLogger(true)

	// Define input and expected output
	input := getMockGetMediaInput()
	movies := getMockGetMediaOutput()

	// Create mock GetMediaUsecase
	mockCtrl := gomock.NewController(t)
	mockGetMediaUsecase := NewMockGetMediaUseCase(mockCtrl)
	mockGetMediaUsecase.EXPECT().
		Execute(gomock.Any(), gomock.Eq(input)).
		Return(movies, nil).
		Times(1)

	// Initialize server with the mock usecase
	handlers := &Handlers{
		Logger:          logger,
		GetMediaUseCase: mockGetMediaUsecase,
	}
	router := rtr.InitRouter(handlers, logger, "/")
	server := srv.NewServer(router)

	// Create a New Request
	requestURL := fmt.Sprintf("/media/%d", input.MediaID)
	req, err := http.NewRequest("GET", requestURL, nil)
	require.NoError(t, err)

	// Execute Request
	response := executeRequest(logger, server, req)

	// Assert the response
	checkResponse(t, logger, response, http.StatusOK, movies)
}

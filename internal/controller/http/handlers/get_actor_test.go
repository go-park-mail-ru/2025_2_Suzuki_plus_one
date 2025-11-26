package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	rtr "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/router"
	srv "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http"
	. "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/handlers"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/stretchr/testify/require"

	. "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/mocks"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	gomock "go.uber.org/mock/gomock"
)

func getMockGetActorInput() dto.GetActorInput {
	return dto.GetActorInput{
		ActorID: 2,
	}
}

func getMockGetActorOutput() dto.GetActorOutput {
	return dto.GetActorOutput{}
}

// Call GetActor handler and check response with query parameters limit and offset
func TestGetActor(t *testing.T) {
	logger := logger.NewZapLogger(true)

	// Define input and expected output
	input := getMockGetActorInput()
	movies := getMockGetActorOutput()

	// Create mock GetActorUsecase
	mockCtrl := gomock.NewController(t)
	mockGetActorUsecase := NewMockGetActorUseCase(mockCtrl)
	mockGetActorUsecase.EXPECT().
		Execute(gomock.Any(), gomock.Eq(input)).
		Return(movies, nil).
		Times(1)

	// Initialize server with the mock usecase
	handlers := &Handlers{
		Logger:          logger,
		GetActorUseCase: mockGetActorUsecase,
	}
	router := rtr.InitRouter(handlers, logger, "/")
	server := srv.NewServer(router)

	// Create a New Request
	requestURL := fmt.Sprintf("/actor/%d", input.ActorID)
	req, err := http.NewRequest("GET", requestURL, nil)
	require.NoError(t, err)

	// Execute Request
	response := executeRequest(logger, server, req)

	// Assert the response
	checkResponse(t, logger, response, http.StatusOK, movies)
}

package handlers_test

import (
	"fmt"
	"net/http"
	"strconv"
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

func getMockGetMediaRecommendationsInput() dto.GetMediaRecommendationsInput {
	return dto.GetMediaRecommendationsInput{
		Limit:  2,
		Offset: 3,
		Type:   "movie",
	}
}

func getMockGetMediaRecommendationsOutput() dto.GetMediaRecommendationsOutput {
	// Create 5 mock movies
	movies := []dto.GetMediaOutput{}
	for i := 0; i < 5; i++ {
		// Create new movie with ID and Title
		movie := dto.GetMediaOutput{}
		movie.MediaID = uint(i + 1)
		movie.Title = "Movie " + strconv.Itoa(i+1)
		movies = append(movies, movie)
	}

	return dto.GetMediaRecommendationsOutput{
		Movies: movies,
	}
}

// Call GetMediaRecommendations handler and check response with query parameters limit and offset
func TestGetMediaRecommendations(t *testing.T) {
	logger := logger.NewZapLogger(true)

	// Define input and expected output
	input := getMockGetMediaRecommendationsInput()
	movies := getMockGetMediaRecommendationsOutput()

	// Create mock GetMediaRecommendationsUsecase
	mockCtrl := gomock.NewController(t)
	mockGetMediaRecommendationsUsecase := NewMockGetMediaRecommendationsUseCase(mockCtrl)
	mockGetMediaRecommendationsUsecase.EXPECT().
		Execute(gomock.Any(), gomock.Eq(input)).
		Return(movies, nil).
		Times(1)

	// Initialize server with the mock usecase
	handlers := &Handlers{
		Logger:                         logger,
		GetMediaRecommendationsUseCase: mockGetMediaRecommendationsUsecase,
	}
	router := rtr.InitRouter(handlers, logger, "/")
	server := srv.NewServer(router)

	// Create a New Request
	requestURL := fmt.Sprintf("/media/recommendations?limit=%d&offset=%d&type=%s", input.Limit, input.Offset, input.Type)
	req, err := http.NewRequest("GET", requestURL, nil)
	require.NoError(t, err)

	// Execute Request
	response := executeRequest(logger, server, req)

	// Assert the response
	checkResponse(t, logger, response, http.StatusOK, movies)
}

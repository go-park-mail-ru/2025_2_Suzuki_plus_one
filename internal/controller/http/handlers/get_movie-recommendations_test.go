package handlers_test

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	srv "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http"
	. "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/handlers"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/stretchr/testify/require"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	gomock "go.uber.org/mock/gomock"
)

func getMockGetMovieRecommendationsInput() dto.GetMovieRecommendationsInput {
	return dto.GetMovieRecommendationsInput{
		Limit:  2,
		Offset: 3,
	}
}

func getMockGetMovieRecommendationsOutput() dto.GetMovieRecommendationsOutput {
	// Create 5 mock movies
	movies := []dto.Movie{}
	for i := 0; i < 5; i++ {
		// Create new movie with ID and Title
		movie := dto.Movie{}
		movie.ID = i + 1
		movie.Title = "Movie " + strconv.Itoa(i+1)
		movies = append(movies, movie)
	}

	return dto.GetMovieRecommendationsOutput{
		Movies: movies,
	}
}

// Call GetMovieRecommendations handler and check response with query parameters limit and offset
func TestGetMovieRecommendations(t *testing.T) {
	logger := logger.NewZapLogger(true)

	// Define input and expected output
	input := getMockGetMovieRecommendationsInput()
	movies := getMockGetMovieRecommendationsOutput()

	// Create mock GetMovieRecommendationsUsecase
	mockCtrl := gomock.NewController(t)
	mockGetMovieRecommendationsUsecase := controller.NewMockGetMovieRecommendationsUsecase(mockCtrl)
	mockGetMovieRecommendationsUsecase.EXPECT().
		Execute(gomock.Any(), gomock.Eq(input)).
		Return(movies, nil).
		Times(1)

	// Initialize server with the mock usecase
	handlers := NewHandlers(logger, mockGetMovieRecommendationsUsecase, nil)
	router := srv.InitRouter(handlers, logger, "/")
	server := srv.NewServer(router)

	// Create a New Request
	requestURL := fmt.Sprintf("/movie/recommendations?limit=%d&offset=%d", input.Limit, input.Offset)
	req, err := http.NewRequest("GET", requestURL, nil)
	require.NoError(t, err)

	// Execute Request
	response := executeRequest(logger, server, req)

	// Assert the response
	checkResponse(t, logger, response, http.StatusOK, movies)
}

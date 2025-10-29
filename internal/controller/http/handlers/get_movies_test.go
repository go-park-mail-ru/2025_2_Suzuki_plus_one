package handlers_test

import (
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

func getMockGetMoviesInput() dto.GetMoviesInput {
	return dto.GetMoviesInput{
		Limit:  2,
		Offset: 3,
	}
}

func getMockGetMoviesOutput() dto.GetMoviesOutput {
	// Create 5 mock movies
	movies := []dto.Movie{}
	for i := 0; i < 5; i++ {
		movie := dto.NewTestMovie("movie_id_" + strconv.Itoa(i))
		movies = append(movies, movie)
	}

	return dto.GetMoviesOutput{
		Movies: movies,
	}
}

func TestGetMovies(t *testing.T) {
	logger := logger.NewZapLogger(true)

	// Define input and expected output
	input := getMockGetMoviesInput()
	movies := getMockGetMoviesOutput()

	// Create mock GetMoviesUsecase
	mockCtrl := gomock.NewController(t)
	mockGetMoviesUsecase := controller.NewMockGetMoviesUsecase(mockCtrl)
	mockGetMoviesUsecase.EXPECT().
		Execute(gomock.Any(), gomock.Eq(input)).
		Return(movies, nil).
		Times(1)

	// Initialize server with the mock usecase
	handlers := NewHandlers(mockGetMoviesUsecase, logger)
	router := srv.InitRouter(handlers, logger, "/")
	server := srv.NewServer(router)

	// Create a New Request
	req, err := http.NewRequest("GET", "/movies?limit=2&offset=3", nil)
	require.NoError(t, err)

	// Execute Request
	response := executeRequest(logger, server, req)

	// Assert the response
	checkResponse(t, logger, response, http.StatusOK, movies)
}

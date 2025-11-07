package handlers_test

import (
	"fmt"
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

func getMockGetObjectInput() dto.GetObjectInput {
	return dto.GetObjectInput{
		Key:        "1",
		BucketName: "posters",
	}
}

func getMockGetObjectOutput() dto.GetObjectOutput {
	return dto.GetObjectOutput{
		URL: "https://example.com/posters/1",
	}
}

// Call GetObjectURL handler and check response with query parameters limit and offset
func TestGetObject(t *testing.T) {
	logger := logger.NewZapLogger(true)

	// Define input and expected output
	input := getMockGetObjectInput()
	output := getMockGetObjectOutput()

	// Create mock GetObjectUsecase
	mockCtrl := gomock.NewController(t)
	mockGetObjectUsecase := NewMockGetObjectUseCase(mockCtrl)
	mockGetObjectUsecase.EXPECT().
		Execute(gomock.Any(), gomock.Eq(input)).
		Return(output, nil).
		Times(1)

	// Initialize server with the mock usecase
	handlers := &Handlers{
		Logger:                logger,
		GetObjectMediaUseCase: mockGetObjectUsecase,
	}
	router := srv.InitRouter(handlers, logger, "/")
	server := srv.NewServer(router)

	// Create a New Request
	requestURL := fmt.Sprintf("/object?key=%s&bucket_name=%s", input.Key, input.BucketName)
	req, err := http.NewRequest("GET", requestURL, nil)
	require.NoError(t, err)

	// Execute Request
	response := executeRequest(logger, server, req)

	// Assert the response
	checkResponse(t, logger, response, http.StatusOK, output)
}

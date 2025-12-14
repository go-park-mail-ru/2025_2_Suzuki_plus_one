package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	srv "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http"
	. "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/handlers"
	. "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/mocks"
	rtr "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http/router"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/stretchr/testify/require"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	gomock "go.uber.org/mock/gomock"
)

func getMockPostAuthSignInInput() dto.PostAuthSignInInput {
	return dto.PostAuthSignInInput{
		Email:    "test@example.com",
		Password: "password123",
	}
}

func getMockPostAuthSignInOutput() dto.PostAuthSignInOutput {
	return dto.PostAuthSignInOutput{
		AccessToken:  "accessTokenValue",
		RefreshToken: "refreshTokenValue",
	}
}

// Call PostAuthSignIn handler and check response
func TestPostAuthSignIn(t *testing.T) {
	logger := logger.NewZapLogger(true)

	// Define input and expected output
	input := getMockPostAuthSignInInput()
	output := getMockPostAuthSignInOutput()

	// Create mock PostAuthSignInUsecase
	mockCtrl := gomock.NewController(t)
	mockPostAuthSignInUsecase := NewMockPostAuthSignInUseCase(mockCtrl)
	mockPostAuthSignInUsecase.EXPECT().
		Execute(gomock.Any(), gomock.Eq(input)).
		Return(output, nil).
		Times(1)

	// Initialize server with the mock usecase
	handlers := &Handlers{
		Logger:                logger,
		PostAuthSignInUseCase: mockPostAuthSignInUsecase,
	}
	router := rtr.InitRouter(handlers, logger, "/")
	server := srv.NewServer(router)

	// Marshal input to JSON
	body, err := json.Marshal(input)
	require.NoError(t, err)
	reader := io.NopCloser(bytes.NewReader(body))

	// Create request
	req, err := http.NewRequest("POST", "/auth/signin", reader)
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)

	// Execute Request
	response := executeRequest(logger, server, req)

	// Assert the response
	checkResponse(t, logger, response, http.StatusOK, UpdatePostAuthSignInOutput(output))
}

package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/stretchr/testify/require"
)

// executeRequest, creates a new ResponseRecorder
// then executes the request by calling ServeHTTP in the router
// after which the handler writes the response to the response recorder
// which we can then inspect.
func executeRequest(logger logger.Logger, s http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.ServeHTTP(rr, req)
	logger.Info("Request executed",
		logger.ToString("method", req.Method),
		logger.ToString("path", req.URL.Path),
		logger.ToInt("response_code", rr.Code),
	)

	return rr
}

// checkResponseCode is a simple utility to assert the response code and body
func checkResponse(t *testing.T, logger logger.Logger, response *httptest.ResponseRecorder, expectedCode int, expectedBody dto.DTO) {
	logger.Debug("Checking response: ", "response", response)

	// Check the response code
	require.Equal(t, expectedCode, response.Code, "Wrong response code")

	// Check the response body
	require.NotNil(t, response.Body, "Response body is nil")

	// Read and unmarshal the response body
	var actualMovies dto.GetMoviesOutput
	err := json.NewDecoder(response.Body).Decode(&actualMovies)
	require.NoError(t, err, "Error unmarshaling response body")

	// Assert the response body
	require.Equal(t, expectedBody, actualMovies, "Response body is not as expected")
}

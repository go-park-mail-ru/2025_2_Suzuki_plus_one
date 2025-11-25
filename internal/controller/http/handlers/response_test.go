package handlers

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test структуры
func TestResponse(t *testing.T) {
	response := Response{
		Code: http.StatusOK,
		Data: "test data",
	}

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "test data", response.Data)
}

func TestResponseError(t *testing.T) {
	err := errors.New("test error")
	responseErr := ResponseError{
		Code:    http.StatusBadRequest,
		Message: err,
		Details: "test details",
	}

	assert.Equal(t, http.StatusBadRequest, responseErr.Code)
	assert.Equal(t, err, responseErr.Message)
	assert.Equal(t, "test details", responseErr.Details)
}

package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/auth"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: Refactor and fix this tests

const SECRET = "secret"

func testAuthMiddleware(t *testing.T, token string, expectNoTokenFound, expectTokenFailed bool, expectClaims auth.JWTClaims) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	if token != "" {
		req.Header.Set("Authorization", token)
	}

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		// Extract authContext from request context
		val := r.Context().Value(AuthContextKey)
		authCtx, ok := val.(*authContext)
		require.True(t, ok)

		// Check authContext fields
		assert.Equal(t, expectNoTokenFound, authCtx.NoTokenFound)
		assert.Equal(t, expectTokenFailed, authCtx.TokenFailed)
		assert.Equal(t, expectClaims, authCtx.Claims)
	})

	handler := authMiddleware(next, SECRET)
	handler.ServeHTTP(rr, req)
	assert.True(t, called)
}

func TestAuthMiddleware_AllCases(t *testing.T) {
	testAuthMiddleware(t, "", true, false, auth.JWTClaims{})
}

func TestAuthMiddleware_InvalidTokenFormat(t *testing.T) {
	testAuthMiddleware(t, "badtoken", true, false, auth.JWTClaims{})
}

func TestWithAuthRequired_NoToken(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), AuthContextKey, &authContext{NoTokenFound: true})
	req = req.WithContext(ctx)

	handler := http.HandlerFunc(withAuthRequired(func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not call next handler")
	}))

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	var resp models.ErrorResponse
	_ = json.NewDecoder(bytes.NewReader(rr.Body.Bytes())).Decode(&resp)
	assert.Equal(t, ErrNoTokenProvided.Message, resp.Message)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestWithAuthRequired_TokenFailed(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), AuthContextKey, &authContext{TokenFailed: true})
	req = req.WithContext(ctx)

	handler := http.HandlerFunc(withAuthRequired(func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not call next handler")
	}))

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	var resp models.ErrorResponse
	_ = json.NewDecoder(bytes.NewReader(rr.Body.Bytes())).Decode(&resp)
	assert.Equal(t, ErrInvalidOrExpired.Message, resp.Message)
}

// Inject valid token into authContext
func TestWithAuthRequired_ValidInject(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	claims := auth.JWTClaims{}
	ctx := context.WithValue(req.Context(), AuthContextKey, &authContext{
		Claims:       claims,
		NoTokenFound: false,
		TokenFailed:  false,
	})
	req = req.WithContext(ctx)

	called := false
	handler := http.HandlerFunc(withAuthRequired(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)
	assert.True(t, called)
	assert.Equal(t, http.StatusOK, rr.Code)
}

// Add valid token in Authorization header
// To get 200 withAuthRequired has to be prepended with authMiddleware but it is not, then got 401
func TestWithAuthRequired_ValidHeader(t *testing.T) {
	claims := auth.JWTClaims{
		// Fill with required fields if needed, e.g. UserID, Exp, etc.
		Type: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "testApp",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Subject:   "user123",
		},
	}
	token, err := auth.NewAuth("secret").GenerateToken(&claims)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", token)

	called := false
	handler := http.HandlerFunc(withAuthRequired(func(w http.ResponseWriter, r *http.Request) {
		called = true
		require.Fail(t, "should not reach final handler")
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)
	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

// With full chain: logging -> authMiddleware -> withAuthRequired -> final
func TestAuthChain_InvalidToken(t *testing.T) {
	// Prepare server
	server := http.NewServeMux()
	called := false
	next := func(w http.ResponseWriter, r *http.Request) {
		called = true
		require.Fail(t, "should not reach final handler")
	}
	server.HandleFunc("/", withAuthRequired(next))

	finalHandler := loggingMiddleware(
		authMiddleware(
			server,
			"secret",
		),
	)

	// Send request with invalid token
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "badtoken")

	// Call final handler
	finalHandler.ServeHTTP(rr, req)

	assert.False(t, called) // Should not reach final handler
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	var resp models.ErrorResponse
	_ = json.NewDecoder(bytes.NewReader(rr.Body.Bytes())).Decode(&resp)
	assert.Equal(t, ErrInvalidOrExpired.Message, resp.Message)
}

// With full chain: logging -> authMiddleware -> withAuthRequired -> final
func TestAuthChain_ValidToken(t *testing.T) {
	// Generate a valid token using the real JWT logic
	claims := auth.JWTClaims{
		Type: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "testApp",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Subject:   "user123",
		},
	}
	token, err := auth.NewAuth(SECRET).GenerateToken(&claims)
	require.NoError(t, err)

	// Prepare server
	server := http.NewServeMux()
	called := false
	next := func(w http.ResponseWriter, r *http.Request) {
		called = true
		val := r.Context().Value(AuthContextKey)
		require.NotNil(t, val)
		authCtx, ok := val.(*authContext)
		require.NotNil(t, authCtx)
		require.True(t, ok)
		assert.False(t, authCtx.NoTokenFound)
		assert.Equal(t, claims.RegisteredClaims.Subject, authCtx.Claims.Subject)
		assert.False(t, authCtx.TokenFailed)
		w.WriteHeader(http.StatusOK)
	}

	server.HandleFunc("/", withAuthRequired(next))

	finalHandler := loggingMiddleware(
		authMiddleware(
			server,
			SECRET,
		),
	)

	// Send request with valid token
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", token)

	// Call final handler
	finalHandler.ServeHTTP(rr, req)

	assert.True(t, called)
	require.Equal(t, http.StatusOK, rr.Code)
}

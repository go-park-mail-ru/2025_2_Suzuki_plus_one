// This file contains tests for middleware flow and request context handling

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

const SECRET = "secret"
const BADTOKEN = "badtoken"

func getValidClaims(t *testing.T) auth.JWTClaims {
	return auth.JWTClaims{
		Type: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "testApp",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Subject:   "user123",
		},
	}
}

// TODO: Move TestResponseModel to models, add New contructors and serealization to other models
// Create interface like ResponseType in testAuthChain

type TestResponseModel struct {
	CheapLiquorOnIce string `json:"msg"`
}

func getValidResponseModel() TestResponseModel {
	return TestResponseModel{
		CheapLiquorOnIce: "nice",
	}
}

// Returns default handler that requires given claims to match request context.
// Writes 200 or fails the test.
// Needs to be wrapped with [http.HandlerFunc] (also via [withAuthRequired])
func getValidHandlerFunc(t *testing.T, claims auth.JWTClaims) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract authContext from request context
		val := r.Context().Value(AuthContextKey)
		require.NotNil(t, val)
		authCtx, ok := val.(*authContext)
		require.NotNil(t, authCtx)
		require.True(t, ok)

		// Check authContext fields to be ok
		assert.False(t, authCtx.NoTokenFound)
		require.Equal(t, claims, authCtx.Claims)
		assert.False(t, authCtx.TokenFailed)

		// All good, write 200
		w.WriteHeader(http.StatusOK)

		// Send valid response
		resp := getValidResponseModel()
		err := json.NewEncoder(w).Encode(resp)
		require.NoError(t, err)
	}
}

// Do not set token at all
func setNoToken(t *testing.T, req *http.Request) {
	// Do nothing
}

// Set invalid token
func setInvalidToken(t *testing.T, req *http.Request) {
	authenticator := auth.NewAuth(SECRET)
	authenticator.TokenMgr.Set(req, BADTOKEN)
}

// Set token using auth logic
func setValidToken(t *testing.T, req *http.Request, claims auth.JWTClaims) {
	authenticator := auth.NewAuth(SECRET)
	token, err := authenticator.GenerateToken(&claims)
	require.NoError(t, err)
	authenticator.TokenMgr.Set(req, token)
}

// Test authMiddleware directly with given token in Authorization header
// Check that context claims are set properly
func testAuthMiddleware(
	t *testing.T,
	tokenSetter func(*testing.T, *http.Request),
	expectNoTokenFound,
	expectTokenFailed bool,
	expectClaims auth.JWTClaims,
) {

	// Prepare request and response recorder
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	tokenSetter(t, req)

	called := false
	// ? Handler can be refactored with getBasicResponseFunc
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		// Extract authContext from request context
		val := r.Context().Value(AuthContextKey)
		authCtx, ok := val.(*authContext)
		require.True(t, ok)

		// Check authContext fields
		assert.Equal(t, expectNoTokenFound, authCtx.NoTokenFound,
			"expectNoTokenFound=%v, got %v", expectNoTokenFound, authCtx.NoTokenFound)
		assert.Equal(t, expectTokenFailed, authCtx.TokenFailed,
			"expectTokenFailed=%v, got %v", expectTokenFailed, authCtx.TokenFailed)
		assert.Equal(t, expectClaims, authCtx.Claims,
			"expectClaims=%v, got %v", expectClaims, authCtx.Claims)
	})

	handler := authMiddleware(next, SECRET)
	handler.ServeHTTP(rr, req)
	assert.True(t, called)
}

// Expect no token found
func TestAuthMiddleware_BlankToken(t *testing.T) {
	testAuthMiddleware(t, setNoToken, true, false, auth.JWTClaims{})
}

// Expect token to be invalid
func TestAuthMiddleware_InvalidTokenFormat(t *testing.T) {
	testAuthMiddleware(t, setInvalidToken, false, true, auth.JWTClaims{})
}

// Expect token to be valid
func TestAuthMiddleware_ValidToken(t *testing.T) {
	// Prepare claims for token
	claims := getValidClaims(t)
	authenticator := auth.NewAuth(SECRET)
	token, err := authenticator.GenerateToken(&claims)
	require.NoError(t, err)

	// Function that sets valid token in request header
	tokenSetter := func(t *testing.T, req *http.Request) {
		authenticator.TokenMgr.Set(req, token)
	}

	testAuthMiddleware(t, tokenSetter, false, false, claims)
}

// Call withAuthRequired with params injected into context
func testWithAuthRequired(t *testing.T, noTokenFound, tokenFailed bool, claims auth.JWTClaims) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), AuthContextKey, &authContext{
		NoTokenFound: noTokenFound,
		TokenFailed:  tokenFailed,
		Claims:       claims,
	})
	req = req.WithContext(ctx)

	// Call response handler with withAuthRequired middleware
	http.HandlerFunc(withAuthRequired(getValidHandlerFunc(t, claims))).ServeHTTP(rr, req)

	if noTokenFound || tokenFailed {
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		var resp models.ErrorResponse
		err := json.NewDecoder(bytes.NewReader(rr.Body.Bytes())).Decode(&resp)
		require.NoError(t, err)
		if noTokenFound {
			assert.Equal(t, ErrNoTokenProvided.Message, resp.Message)
		} else {
			assert.Equal(t, ErrInvalidOrExpired.Message, resp.Message)
		}
	} else {
		assert.Equal(t, http.StatusOK, rr.Code)
	}
}

func TestWithAuthRequired_NoToken(t *testing.T) {
	testWithAuthRequired(t, true, false, auth.JWTClaims{})
}

func TestWithAuthRequired_TokenFailed(t *testing.T) {
	testWithAuthRequired(t, false, true, auth.JWTClaims{})
}

func TestWithAuthRequired_ValidToken(t *testing.T) {
	testWithAuthRequired(t, false, false, getValidClaims(t))
}

// Add valid token in Authorization header, not in the context
// To get 200 withAuthRequired has to be prepended with authMiddleware but it is not, then got 401
func TestWithAuthRequired_ValidHeader(t *testing.T) {
	// Prepare request and response recorder
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	// Set valid token in header
	claims := getValidClaims(t)
	authenticator := auth.NewAuth(SECRET)
	token, err := authenticator.GenerateToken(&claims)
	require.NoError(t, err)
	authenticator.TokenMgr.Set(req, token)

	handler := func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Can't reach handler without authContext in request")
	}
	http.HandlerFunc(withAuthRequired(handler)).ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

// Test request as a full request chain:
//
//	logging -> authMiddleware -> withAuthRequired -> handler
//
// Uses [getValidResponseFunc] as a request handler
// Then compares expectedResponse with response body
func testAuthChain[ResponseType any](t *testing.T,
	tokenSetter func(*testing.T, *http.Request),
	expectedClaims auth.JWTClaims, expectedStatus int, expectedResponse ResponseType) {

	// Prepare server
	server := http.NewServeMux()
	server.HandleFunc("/", withAuthRequired(getValidHandlerFunc(t, expectedClaims)))

	// Set middleware chain
	finalHandler := loggingMiddleware(
		authMiddleware(
			server,
			SECRET,
		),
	)

	// Send request with invalid token
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	// Set token in header
	tokenSetter(t, req)

	// Call final handler
	finalHandler.ServeHTTP(rr, req)

	// Check status code to be as expected
	assert.Equal(t, expectedStatus, rr.Code)

	// Check response to be as expected
	var resp ResponseType
	err := json.NewDecoder(bytes.NewReader(rr.Body.Bytes())).Decode(&resp)
	if err != nil {
		t.Fatalf(
			"Failed to decode expected response property: err=%v, expectedResponse=%+v",
			err,
			expectedResponse,
		)
	}
	assert.Equal(t, expectedResponse, resp)
}

// Pass no token on [getValidHandlerFunc], expect 401 from withAuthRequired
func TestAuthChain_NoToken(t *testing.T) {
	tokenSetter := setNoToken
	claims := auth.JWTClaims{}
	expectedStatus := http.StatusUnauthorized
	expectedResponse := ErrNoTokenProvided
	testAuthChain(t, tokenSetter, claims, expectedStatus, expectedResponse)
}

// Pass invalid token on [getValidHandlerFunc], expect 401 from withAuthRequired
func TestAuthChain_InvalidToken(t *testing.T) {
	tokenSetter := setInvalidToken
	claims := auth.JWTClaims{}
	expectedStatus := http.StatusUnauthorized
	expectedResponse := ErrInvalidOrExpired
	testAuthChain(t, tokenSetter, claims, expectedStatus, expectedResponse)
}

// Pass valid token on withAuthRequired,
// Expect 200 response from [getValidResponseFunc] (bc testAuthChain uses it)
func TestAuthChain_ValidToken(t *testing.T) {
	claims := getValidClaims(t)
	tokenSetter := func(t *testing.T, req *http.Request) {
		setValidToken(t, req, claims)
	}
	expectedStatus := http.StatusOK
	expectedResponse := getValidResponseModel()
	testAuthChain(t, tokenSetter, claims, expectedStatus, expectedResponse)
}

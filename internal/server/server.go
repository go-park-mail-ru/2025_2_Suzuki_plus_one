package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/auth"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/config"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/db"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/models"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/utils"
)

type Server struct {
	address           string         // Serving address
	authSecret        string         // JWT secret for HMAC
	serverName        string         // Server name
	accessTokenExpiry time.Duration  // JWT access token expiry duration
	frontendOrigin    string         // Frontend origin for CORS
	db                *db.DataBase   // Database connection
	server            *http.ServeMux // HTTP request multiplexer
}

func NewServer(cfg *config.Config, database *db.DataBase) *Server {
	mux := http.NewServeMux()
	return &Server{
		address:           cfg.SERVER_SERVE_STRING,
		authSecret:        cfg.SERVER_JWT_SECRET,
		serverName:        cfg.SERVER_NAME,
		accessTokenExpiry: cfg.SERVER_JWT_ACCESS_EXPIRATION,
		frontendOrigin:    cfg.SERVER_FRONTEND_URL,
		db:                database,
		server:            mux,
	}
}

// Get all movies from database
func (s *Server) getAllMovies(w http.ResponseWriter, r *http.Request) {
	request := models.MoviesRequest{}

	// Parse query parameters
	query := r.URL.Query()
	if offStr := query.Get("offset"); offStr != "" {
		// If it is black we leave it as default 0
		if _, err := fmt.Sscanf(offStr, "%d", &request.Offset); err != nil {
			responseWithError(w, http.StatusBadRequest, ErrMoviesInvalidParams)
			return
		}
	}
	if limStr := query.Get("limit"); limStr != "" {
		if _, err := fmt.Sscanf(limStr, "%d", &request.Limit); err != nil {
			responseWithError(w, http.StatusBadRequest, ErrMoviesInvalidParams)
			return
		}
	}

	movies := s.db.GetMovies(request.Offset, request.Limit)
	log.Printf("Fetched %d movies from database", len(movies))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}

func (s *Server) signUp(w http.ResponseWriter, r *http.Request) {
	request := models.SignUpRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responseWithError(w, http.StatusBadRequest, ErrSignUpWrongData)
		return
	}

	if request.Email == "" || request.Password == "" {
		log.Println("Email or password is empty")
		responseWithError(w, http.StatusBadRequest, ErrSignUpWrongData)
		return
	}

	err := s.db.CreateUser(request.Email, request.Password)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			responseWithError(w, http.StatusConflict, ErrSignUpUserExists)
		} else {
			responseWithError(w, http.StatusInternalServerError, errorWithDetails(ErrSignUpInternal, err.Error()))
		}
		return
	}

	user := s.db.FindUserByEmail(request.Email)
	if user == nil {
		responseWithError(w, http.StatusInternalServerError, ErrSignUpInternal)
		return
	}

	authenticator := auth.NewAuth(s.authSecret)
	claims := auth.NewJWTClaims(user.ID, "access", s.accessTokenExpiry, s.serverName)
	token, err := authenticator.GenerateToken(claims)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, errorWithDetails(ErrSignUpInternal, err.Error()))
		return
	}

	response := models.SignUpResponse{
		Token: token,
		User: models.UserAPI{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) signIn(w http.ResponseWriter, r *http.Request) {
	request := models.SignInRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responseWithError(w, http.StatusBadRequest, ErrSignInWrongData)
		return
	}

	// Query user from database
	user := s.db.FindUserByEmail(request.Email)
	if user == nil {
		responseWithError(w, http.StatusUnauthorized, ErrSignInWrongData)
		return
	}

	// Check password
	if err := utils.ValidateHashedPasswordBcrypt(user.PasswordHash, request.Password); err != nil {
		responseWithError(w, http.StatusUnauthorized, ErrSignInWrongData)
		return
	}

	// Right credentials, create token
	authenticator := auth.NewAuth(s.authSecret)
	claims := auth.NewJWTClaims(user.ID, "access", s.accessTokenExpiry, s.serverName)
	token, err := authenticator.GenerateToken(claims)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, errorWithDetails(ErrSignInInternal, err.Error()))
		return
	}

	response := models.SignInResponse{
		Token: token,
		User: models.UserAPI{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) signOut(w http.ResponseWriter, r *http.Request) {
	response := models.SignOutResponse{
		Success: true,
		Message: "Successfully signed out",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) setupRoutes() {
	s.server.HandleFunc("/movies", s.getAllMovies)
	s.server.HandleFunc("/auth/signup", s.signUp)
	s.server.HandleFunc("/auth/signin", s.signIn)
	s.server.Handle("/auth/signout", withAuthRequired(http.HandlerFunc(s.signOut)))
}

func (s *Server) Serve() {
	s.setupRoutes()

	// Add middleware
	logHandler := loggingMiddleware(s.server)
	jsonHandler := forceJSONMiddleware(logHandler)
	corsHandler := corsMiddleware(jsonHandler, s.frontendOrigin)
	authedHandler := authMiddleware(corsHandler, s.authSecret)
	handler := authedHandler

	log.Println("Server starting on", s.address)
	log.Fatal(http.ListenAndServe(s.address, handler))
}

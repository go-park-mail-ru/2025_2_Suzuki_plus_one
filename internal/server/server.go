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
	prefix            string         // API prefix
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
		prefix:            cfg.SERVER_SERVE_PREFIX,
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
		// If parameter is blank we leave it as default 0
		if _, err := fmt.Sscanf(offStr, "%d", &request.Offset); err != nil {
			responseWithError(w, http.StatusBadRequest, ErrMoviesInvalidParams)
			return
		}
	}
	if limStr := query.Get("limit"); limStr != "" {
		// If parameter is blank we leave it as default 0 (means no limit)
		if _, err := fmt.Sscanf(limStr, "%d", &request.Limit); err != nil {
			responseWithError(w, http.StatusBadRequest, ErrMoviesInvalidParams)
			return
		}
	}

	movies := s.db.FindMovies(request.Offset, request.Limit)

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
		// TODO: create a fancy logger
		log.Println("SignUp: Email or password is empty")
		responseWithError(w, http.StatusBadRequest, ErrSignUpWrongData)
		return
	}

	// Create user in database or get existing one
	user, err, created := s.db.CreateUser(request.Email, request.Password)
	if err != nil {
		if !created {
			responseWithError(w, http.StatusConflict, errorWithDetails(ErrSignUpUserExists, err.Error()))
		} else {
			responseWithError(w, http.StatusInternalServerError, errorWithDetails(ErrSignUpInternal, err.Error()))
		}
		return
	}
	log.Println("SignUp: Created new user:", user.Email)

	// Create token for the new user
	authenticator := auth.NewAuth(s.authSecret)
	claims := auth.NewJWTClaims(user.ID, "access", s.accessTokenExpiry, s.serverName)
	token, err := authenticator.GenerateToken(claims)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, errorWithDetails(ErrSignUpInternal, err.Error()))
		return
	}

	// Return signup response
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

	log.Println("SignIn: User signed in:", user.Email)
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
	// Since we use stateless JWT, sign-out is handled on the client side by deleting the token.
	// TODO: Think about token blacklisting
	// Context is already checked by withAuthRequired middleware
	log.Println("SignOut: User signed out")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) auth(w http.ResponseWriter, r *http.Request) {
	// Context is already checked by withAuthRequired middleware
	authCtx := r.Context().Value(AuthContextKey).(*authContext)

	user := s.db.FindUserByID(authCtx.Claims.Subject)
	if user == nil {
		log.Println("Auth: User not found")
		responseWithError(w, http.StatusUnauthorized, ErrInvalidOrExpired)
		return
	}
	log.Println("Auth: User authenticated:", user.Email)
	response := models.AuthResponse{
		User: models.UserAPI{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Add handlers to routes
// Prefix is used for versioning, e.g. /api/v1/
func (s *Server) setupRoutes(prefix string) {
	s.server.HandleFunc(prefix+"/movies", s.getAllMovies)
	s.server.HandleFunc(prefix+"/auth/signup", s.signUp)
	s.server.HandleFunc(prefix+"/auth/signin", s.signIn)
	s.server.HandleFunc(prefix+"/auth/signout", withAuthRequired(s.signOut))
	s.server.HandleFunc(prefix+"/auth", withAuthRequired(s.auth))
}

func (s *Server) Serve() {
	s.setupRoutes(s.prefix)
	
	// Add middleware, the order is important
	// Request -> Logging -> CORS -> Auth -> JSON -> (-> require Auth -> ...) Handlers -> Response
	handler := loggingMiddleware(
		corsMiddleware(
			authMiddleware(
				forceJSONMiddleware(s.server),
				s.authSecret,
			),
			s.frontendOrigin,
		),
	)

	log.Println("Server starting on", s.address)
	log.Fatal(http.ListenAndServe(s.address, handler))
}

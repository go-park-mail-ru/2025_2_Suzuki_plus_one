package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

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
	logger            *zap.Logger
}

func NewServer(cfg *config.Config, database *db.DataBase, logger *zap.Logger) *Server {
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
		logger:            logger,
	}
	// TODO: store auth instance in Server struct (like db connection)
	// Note: Now, we call NewAuth(authSecret) anytime we need auth service
	// However, we can optimize it by storing auth instance in Server struct
	// But it will require concurrent safety if we store stateful data there
}

// Get all movies from database
func (s *Server) getAllMovies(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Fetching movies",
		zap.String("method", r.Method),
		zap.String("query", r.URL.RawQuery))

	request := models.MoviesRequest{}

	// Parse query parameters
	query := r.URL.Query()
	if offStr := query.Get("offset"); offStr != "" {
		// If parameter is blank we leave it as default 0
		if _, err := fmt.Sscanf(offStr, "%d", &request.Offset); err != nil {
			s.logger.Warn("Invalid offset parameter",
				zap.String("offset", offStr),
				zap.Error(err))

			responseWithError(w, http.StatusBadRequest, ErrMoviesInvalidParams, s.logger)
			return
		}
	}
	if limStr := query.Get("limit"); limStr != "" {
		// If parameter is blank we leave it as default 0 (means no limit)
		if _, err := fmt.Sscanf(limStr, "%d", &request.Limit); err != nil {
			s.logger.Warn("Invalid limit parameter",
				zap.String("limit", limStr),
				zap.Error(err))

			responseWithError(w, http.StatusBadRequest, ErrMoviesInvalidParams, s.logger)
			return
		}
	}

	movies := s.db.FindMovies(request.Offset, request.Limit)

	s.logger.Info("Fetching completed successfully",
		zap.String("count", strconv.FormatInt(int64(len(movies)), 10)),
		zap.String("offset", strconv.FormatUint(uint64(request.Offset), 10)),
		zap.String("limit", strconv.FormatUint(uint64(request.Limit), 10)))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}

func (s *Server) signUp(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Sign up request")
	request := models.SignUpRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.logger.Warn("Failed to decode request", zap.Error(err))
		responseWithError(w, http.StatusBadRequest, ErrSignUpWrongData, s.logger)
		return
	}

	if request.Email == "" || request.Password == "" {
		// TODO: create a fancy logger
		s.logger.Warn("Email or Password is empty")
		responseWithError(w, http.StatusBadRequest, ErrSignUpWrongData, s.logger)
		return
	}

	// Create user in database or get existing one
	user, err, created := s.db.CreateUser(request.Email, request.Password, s.logger)
	if err != nil {
		if !created {
			s.logger.Warn("User already exists",
				zap.String("email", request.Email),
				zap.Error(err))

			responseWithError(w, http.StatusConflict, errorWithDetails(ErrSignUpUserExists, err.Error()), s.logger)
		} else {
			s.logger.Warn("Signup internal error",
				zap.String("email", request.Email),
				zap.Error(err))

			responseWithError(w, http.StatusInternalServerError, errorWithDetails(ErrSignUpInternal, err.Error()), s.logger)
		}
		return
	}
	s.logger.Info("User created successfully",
		zap.String("id", user.ID),
		zap.String("email", request.Email))

	// Create token for the new user
	authenticator := auth.NewAuth(s.authSecret, s.logger)
	claims := auth.NewJWTClaims(user.ID, "access", s.accessTokenExpiry, s.serverName)
	token, err := authenticator.GenerateToken(claims)
	if err != nil {
		s.logger.Warn("Failed to generate token",
			zap.String("id", user.ID),
			zap.Error(err))

		responseWithError(w, http.StatusInternalServerError, errorWithDetails(ErrSignUpInternal, err.Error()), s.logger)
		return
	}

	// TODO: move this to internal/models
	// Return signup response
	response := models.SignInResponse{ // SignUp and SignIn responses are the same
		User: models.UserAPI{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	}
	w.WriteHeader(http.StatusCreated)
	authenticator.TokenMgr.ResponseWithAuth(w, token, response)

	s.logger.Info("User signed up successfully",
		zap.String("id", user.ID))
}

func (s *Server) signIn(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Sign in request")
	request := models.SignInRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.logger.Warn("Failed to decode request", zap.Error(err))

		responseWithError(w, http.StatusBadRequest, ErrSignInWrongData, s.logger)
		return
	}

	// Query user from database
	user := s.db.FindUserByEmail(request.Email)
	if user == nil {
		s.logger.Warn("User does not exist",
			zap.String("email", request.Email))

		responseWithError(w, http.StatusUnauthorized, ErrSignInWrongData, s.logger)
		return
	}

	// Check password
	if err := utils.ValidateHashedPasswordBcrypt(user.PasswordHash, request.Password); err != nil {
		s.logger.Warn("Invalid password",
			zap.String("email", request.Email),
			zap.Error(err))

		responseWithError(w, http.StatusUnauthorized, ErrSignInWrongData, s.logger)
		return
	}

	// Right credentials, create token
	authenticator := auth.NewAuth(s.authSecret, s.logger)
	claims := auth.NewJWTClaims(user.ID, "access", s.accessTokenExpiry, s.serverName)
	token, err := authenticator.GenerateToken(claims)
	if err != nil {
		s.logger.Warn("Failed to generate token",
			zap.String("id", user.ID),
			zap.String("email", request.Email))
		responseWithError(w, http.StatusInternalServerError, errorWithDetails(ErrSignInInternal, err.Error()), s.logger)
		return
	}

	log.Println("SignIn: User signed in:", user.Email)
	// TODO: move this to internal/models converter DB -> API Response
	response := models.SignInResponse{
		User: models.UserAPI{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	}
	w.WriteHeader(http.StatusOK)
	authenticator.TokenMgr.ResponseWithAuth(w, token, response)
}

func (s *Server) signOut(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Sign out request")
	// Since we use stateless JWT, sign-out is handled on the client side by deleting the token.
	// Context is already checked by withAuthRequired middleware

	// TODO: Think about token blacklisting for token key
	// TODO: Implement refresh tokens with short-lived access tokens
	authentication := auth.NewAuth(s.authSecret, s.logger)
	w.WriteHeader(http.StatusOK)
	authentication.TokenMgr.ResponseWithDeauth(w)

	s.logger.Info("User signed out successfully")
}

func (s *Server) auth(w http.ResponseWriter, r *http.Request) {
	// Context is already checked by withAuthRequired middleware
	authCtx := r.Context().Value(AuthContextKey).(*authContext)

	user := s.db.FindUserByID(authCtx.Claims.Subject)
	if user == nil {
		s.logger.Warn("Auth user does not find",
			zap.String("user_id", authCtx.Claims.Subject))

		responseWithError(w, http.StatusUnauthorized, ErrInvalidOrExpired, s.logger)
		return
	}

	s.logger.Debug("User authenticated",
		zap.String("user_id", authCtx.Claims.Subject),
		zap.String("email", user.Email))

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
func (s *Server) setupRoutes(prefix string) error {
	s.server.HandleFunc(prefix+"/movies", s.getAllMovies)
	s.server.HandleFunc(prefix+"/auth/signup", s.signUp)
	s.server.HandleFunc(prefix+"/auth/signin", s.signIn)
	s.server.HandleFunc(prefix+"/auth/signout", withAuthRequired(s.signOut, s.logger))
	s.server.HandleFunc(prefix+"/auth", withAuthRequired(s.auth, s.logger))
}

func (s *Server) Serve() error {
	s.setupRoutes(s.prefix)

	// Add middleware, the order is important
	// Request -> Logging -> CORS -> Auth -> JSON -> (-> require Auth -> ...) Handlers -> Response
	handler := loggingMiddleware(
		corsMiddleware(
			authMiddleware(
				forceJSONMiddleware(s.server),
				s.authSecret,
				s.logger,
			),
			s.frontendOrigin,
			s.logger,
		),
		s.logger,
	)

	s.logger.Info("Server started",
		zap.String("address", s.address),
		zap.String("prefix", s.prefix))

	if err := http.ListenAndServe(s.address, handler); err != nil {
		s.logger.Error("Failed to start server", zap.Error(err))
		return err
	}

	return nil
}

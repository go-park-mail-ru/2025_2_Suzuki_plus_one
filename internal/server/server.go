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

// TODO: remove this example handler
func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

// Get all movies from database
func (s *Server) getAllMovies(w http.ResponseWriter, r *http.Request) {
	// Default values
	offset := 0
	limit := 10

	// Parse query parameters
	query := r.URL.Query()
	if offStr := query.Get("offset"); offStr != "" {
		fmt.Sscanf(offStr, "%d", &offset)
	}
	if limStr := query.Get("limit"); limStr != "" {
		fmt.Sscanf(limStr, "%d", &limit)
	}

	movies := s.db.GetMovies(offset, limit)
	log.Printf("Fetched %d movies from database", len(movies))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}

// Note: This is disabled in API for now
// func (s *Server) getMovieByID(w http.ResponseWriter, r *http.Request) {
// 	id := strings.TrimPrefix(r.URL.Path, "/movies/")

// 	//Here we should get movie by id from database by id
// 	movie := []map[string]interface{}{
// 		{
// 			"id":     id,
// 			"title":  "Movies",
// 			"year":   2010,
// 			"genres": []string{"Sci-Fi", "Thriller"},
// 		},
// 	}

// 	json.NewEncoder(w).Encode(movie)
// }

func (s *Server) signUp(w http.ResponseWriter, r *http.Request) {
	/*
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	*/

	var request struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	/*
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	*/

	//Here must be business-logic of creating new user and adding him to db
	response := map[string]interface{}{
		"success": true,
		"message": "User created successfully",
		"token":   "examplejskdflakdsjf",

		"user": map[string]string{
			"id":       "user1",
			"email":    request.Email,
			"username": request.Username,
		},
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) signIn(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

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
	if err := utils.ValidateHashedPassword(user.PasswordHash, request.Password); err != nil {
		responseWithError(w, http.StatusUnauthorized, ErrSignInWrongData)
		return
	}

	// Right credentials, create token
	// token, err := auth.CreateToken(user.ID, s.authSecret)
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
	// Must be some sign out logic
	response := map[string]interface{}{
		"success": true,
		"message": "Successfully signed out",
	}

	json.NewEncoder(w).Encode(response)
}

func (s *Server) setupRoutes() {
	s.server.HandleFunc("/", greet)
	s.server.HandleFunc("/movies", s.getAllMovies)
	s.server.HandleFunc("/auth/signup", s.signUp)
	s.server.HandleFunc("/auth/signin", s.signIn)
	s.server.Handle("/auth/signout", withAuthRequired(http.HandlerFunc(s.signOut)))

	// Note: Old routing logic, kept for reference // TODO: remove later
	// switch {
	// case path == "/" && method == "GET":
	// 	greet(w, r)
	// case path == "/movies" && method == "GET":
	// 	getAllMovies(w, r)
	// case strings.HasPrefix(path, "/movies/") && method == "GET":
	// 	getMovieByID(w, r)
	// case path == "/auth/signup":
	// 	signUp(w, r)
	// case path == "/auth/signin":
	// 	signIn(w, r)
	// case path == "/auth/signout":
	// 	signOut(w, r)
	// default:
	// 	log.Printf("404 NOT FOUND")
	// 	http.NotFound(w, r)
	// }
}

func (s *Server) Serve() {
	s.setupRoutes()

	// Apply middlewares
	logHandler := loggingMiddleware(s.server)
	jsonHandler := forceJSONMiddleware(logHandler)
	corsHandler := corsMiddleware(jsonHandler, s.frontendOrigin)
	authedHandler := authMiddleware(corsHandler, s.authSecret)
	handler := authedHandler

	log.Println("Server starting on", s.address)
	log.Fatal(http.ListenAndServe(s.address, handler))
}

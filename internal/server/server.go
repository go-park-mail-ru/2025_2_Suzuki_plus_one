package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	db "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/db"
)

type Server struct {
	address string
	db      *db.DataBase
	server  *http.ServeMux
}

func NewServer(address string, database *db.DataBase) *Server {
	mux := http.NewServeMux()
	return &Server{
		address: address,
		db:      database,
		server:  mux,
	}
}

// Middleware for handling CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: restrict this in production
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Middleware for logging requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("<<< REQUEST: %s %s === from %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf(">>> REPLY: Completed in %v", time.Since(start))
	})
}

// Middleware to always set Content-Type to application/json
func alwaysJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

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
		Email    string `json:"email"'`
		Password string `json:"password"`
	}

	/*
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
	*/

	//Must be some authentication logic
	response := map[string]interface{}{
		"success": true,
		"message": "Authentication successfull",
		"token":   "exampleasdkjfpa",

		"user": map[string]string{
			"id":       "user1",
			"email":    request.Email,
			"username": "username",
		},
	}

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
	s.server.HandleFunc("/auth/signout", s.signOut)

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
	loggedHandler := loggingMiddleware(s.server)
	jsonHandler := alwaysJSONMiddleware(loggedHandler)
	handler := corsMiddleware(jsonHandler)

	log.Println("Server starting on", s.address)
	log.Fatal(http.ListenAndServe(s.address, handler))
}

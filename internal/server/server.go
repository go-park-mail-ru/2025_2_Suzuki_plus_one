package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func getAllMovies(w http.ResponseWriter, r *http.Request) {
	//Here we should get (maybe top 10) movies from db
	movies := []map[string]interface{}{
		{
			"id":     "1",
			"title":  "Movies",
			"year":   2010,
			"genres": []string{"Sci-Fi", "Thriller"},
		},
		{
			"id":     "2",
			"title":  "Movies",
			"year":   2011,
			"genres": []string{"Sci-Fi"},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovieByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/movies/")

	//Here we should get movie by id from database by id
	movie := []map[string]interface{}{
		{
			"id":     id,
			"title":  "Movies",
			"year":   2010,
			"genres": []string{"Sci-Fi", "Thriller"},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

func signUp(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func signIn(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func signOut(w http.ResponseWriter, r *http.Request) {
	// Must be some sign out logic
	response := map[string]interface{}{
		"success": true,
		"message": "Successfully signed out",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func router(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method

	log.Printf("=== REQUEST: %s %s ===", method, path)

	switch {
	case path == "/" && method == "GET":
		greet(w, r)
	case path == "/movies" && method == "GET":
		getAllMovies(w, r)
	case strings.HasPrefix(path, "/movies/") && method == "GET":
		getMovieByID(w, r)
	case path == "/auth/signup":
		signUp(w, r)
	case path == "/auth/signin":
		signIn(w, r)
	case path == "/auth/signout":
		signOut(w, r)
	default:
		log.Printf("404 NOT FOUND")
		http.NotFound(w, r)
	}
}

func Serve() {
	http.HandleFunc("/", router)

	log.Println("Server starting on 127.0.0.1:8081")
	log.Fatal(http.ListenAndServe("127.0.0.1:8081", nil))
}

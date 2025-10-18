// This file contains data models for the application.
// It has to correspond to database schema and API spec.

package models

// User model in the database
type UserDB struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

// User model in the API
type UserAPI struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// MovieAPI and MovieDB coincide
type Movie struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Year    int      `json:"year"`
	Genres  []string `json:"genres"`
	Preview string   `json:"preview"`
}

type Actor struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	DateOfBirth  string `json:"dateOfBirth"`
	PlaceOfBirth string `json:"placeOfBirth"`
	Biography    string `json:"biography"`
	Preview      string `json:"preview"`
}

type ActorRequest struct {
	ActorID string `json:"actorId"`
}

// This file contains data models for the application.
// It have to correspond to database schema and API spec.

package models

type UserDB struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

type UserAPI struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type Movie struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Year    int      `json:"year"`
	Genres  []string `json:"genres"`
	Preview string   `json:"preview"`
}

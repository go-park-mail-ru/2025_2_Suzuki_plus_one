package db

import (
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/models"
)

type DataBase struct {
	movies []models.Movie
	users  []models.User
}

// Initialize a new in-memory (for now) database
// It is filled with mock users and movies
func NewDataBase() *DataBase {
	return &DataBase{
		movies: models.MockMovies(),
		users:  models.MockUsers(),
	}
}

// Get all movies from the database (for now, returns all mock movies) // TODO
func (db *DataBase) GetMovies(offset int, limit int) []models.Movie {
	return db.movies
}

// Create a new user in the database
func (db *DataBase) CreateUser(user models.User) {
	db.users = append(db.users, user)
}

// Find a user by email
func (db *DataBase) FindUserByEmail(email string) *models.User {
	for _, user := range db.users {
		if user.Email == email {
			return &user
		}
	}
	return nil
}
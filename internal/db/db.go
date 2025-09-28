package db

import (
	"errors"
	"strconv"
	"sync"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/models"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/utils"
)

type DataBase struct {
	movies []models.Movie
	users  []models.UserDB
	mu     sync.RWMutex // Mutex for thread-safe access
}

// Initialize a new in-memory (for now) database
// It is filled with mock users and movies
func NewDataBase() *DataBase {
	return &DataBase{
		movies: models.MockMovies(),
		users:  models.MockUsers(),
	}
}

// Get all movies from the database (for now, returns all mock movies)
// TODO: implement offset and limit
func (db *DataBase) GetMovies(offset int, limit int) []models.Movie {
	result := make(chan []models.Movie)

	go func() {
		defer close(result)
		db.mu.RLock()
		defer db.mu.RUnlock()

		// Simulate fetching movies (for now, return all movies)
		result <- db.movies
	}()
	return <-result
}

// Generate a new user ID (just a simple increment for now)
func (db *DataBase) GetNewID() string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return strconv.Itoa(len(db.users) + 1)
}

// Create a new user asynchronously
func (db *DataBase) CreateUser(email, password string) error {
	// Validate email and password
	if err := ValidateEmail(email); err != nil {
		return err
	}
	if err := ValidatePassword(password); err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	result := make(chan bool)

	if db.FindUserByEmail(email) != nil {
		return errors.New("user with this email already exists")
	}

	go func() {
		defer close(result)
		db.mu.Lock()
		defer db.mu.Unlock()

		id := db.GetNewID()
		// Simulate creating a user
		user := models.UserDB{
			ID:           id,
			Username:     "user" + id,
			Email:        email,
			PasswordHash: hashedPassword,
		}
		db.users = append(db.users, user)
		result <- true
	}()
	<-result // Wait for the goroutine to finish
	return nil
}

// Find a user by email asynchronously
func (db *DataBase) FindUserByEmail(email string) *models.UserDB {
	result := make(chan *models.UserDB)

	go func() {
		defer close(result)
		db.mu.RLock()
		defer db.mu.RUnlock()

		// Simulate finding a user by email
		for _, user := range db.users {
			if user.Email == email {
				result <- &user
				return
			}
		}
		result <- nil // Return nil if no user is found
	}()

	return <-result
}

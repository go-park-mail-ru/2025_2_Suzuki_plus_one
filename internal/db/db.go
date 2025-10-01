package db

import (
	"strconv"
	"sync"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/models"
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

func NewMockDB() *DataBase {
	return &DataBase{
		movies: models.MockMovies(),
		users:  models.MockUsers(),
	}
}

// Generate a new user ID (just a simple increment for now)
// You need to call this function with a lock held
func (db *DataBase) getNewID() string {
	return strconv.Itoa(len(db.users) + 1)
}

package inmemory

import (
	"sync"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
)

type DataBase struct {
	movies []entity.Movie
	users  []entity.User
	mu     sync.RWMutex // Mutex for thread-safe access
}

// Initialize a new in-memory (for now) database
// It is filled with mock users and movies
func NewDataBase() *DataBase {
	return &DataBase{
		movies: []entity.Movie{},
		users:  []entity.User{},
	}
}

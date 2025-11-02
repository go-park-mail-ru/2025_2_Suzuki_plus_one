package inmemory

import (
	"strconv"
	"sync"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type DataBase struct {
	movies []entity.Movie
	users  []entity.User
	mu     sync.RWMutex // Mutex for thread-safe access
	logger logger.Logger
}

// Initialize a new in-memory (for now) database
// It is filled with mock users and movies
func NewDataBase(logger logger.Logger) *DataBase {
	return &DataBase{
		movies: []entity.Movie{},
		users:  []entity.User{},
		logger: logger,
	}
}

func (db *DataBase) Connect() error {
	// Fill with mock data
	db.mu.Lock()
	defer db.mu.Unlock()

	// Mock movies
	for i := 1; i <= 100; i++ {
		movie := entity.Movie{
			ID:    strconv.Itoa(i),
			Title: "Movie " + strconv.Itoa(i),
			Year:  2000 + i%20,
		}
		db.movies = append(db.movies, movie)
	}

	// Mock users
	for i := 1; i <= 10; i++ {
		user := entity.User{
			ID:       strconv.Itoa(i),
			Username: "user" + strconv.Itoa(i),
			Email:    "user" + strconv.Itoa(i) + "@example.com",
		}
		db.users = append(db.users, user)
	}

	db.logger.Info("In-memory database initialized with mock data")
	return nil
}

func (db *DataBase) Close() error {
	// No resources to close in in-memory database
	db.logger.Info("In-memory database closed")
	return nil
}

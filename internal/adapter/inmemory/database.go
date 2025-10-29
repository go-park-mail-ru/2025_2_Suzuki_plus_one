package inmemory

import (
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

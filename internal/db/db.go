package db

import (
	"strconv"
	"sync"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/models"
)

type DataBase struct {
	movies []models.Movie
	users  []models.UserDB
	actors []models.Actor
	mu     sync.RWMutex // Mutex for thread-safe access
}

// Initialize a new in-memory (for now) database
// It is filled with mock users and movies
func NewDataBase() *DataBase {
	return &DataBase{
		movies: models.MockMovies(),
		users:  models.MockUsers(),
		actors: models.MockActors(),
	}
}

func NewMockDB() *DataBase {
	return &DataBase{
		movies: models.MockMovies(),
		users:  models.MockUsers(),
		actors: models.MockActors(),
	}
}

// Generate a new user ID (just a simple increment for now)
// You need to call this function with a lock held
func (db *DataBase) getNewID() string {
	return strconv.Itoa(len(db.users) + 1)
}

// All these functions must be completed when postgres will be connected
func (db *DataBase) UpdateUserAccount(id string, request models.UserAccountUpdate) error {
	// Just doin nothing rn
	return nil
}

func (db *DataBase) FindUserNotifications(id string) interface{} {
	// Just doin nothing rn
	return nil
}

func (db *DataBase) UpdateUserNotifications(id string, request models.NotificationSettings) error {
	// Just doin nothing rn
	return nil
}

func (db *DataBase) UpdateUserPassword(id string, password string, password2 string) error {
	// Just doin nothing rn
	return nil
}

package db

import (
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/models"
)

// Find a user by email asynchronously
// Returns nil if no user is found
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

func (db *DataBase) FindUserByID(id string) *models.UserDB {
	result := make(chan *models.UserDB)

	go func() {
		defer close(result)
		db.mu.RLock()
		defer db.mu.RUnlock()

		// Simulate finding a user by ID
		for _, user := range db.users {
			if user.ID == id {
				result <- &user
				return
			}
		}
		result <- nil // Return nil if no user is found
	}()

	return <-result
}

// Get movies from the database
// If limit is 0, return all movies from offset
func (db *DataBase) FindMovies(offset uint, limit uint) []models.Movie {
	if offset > uint(len(db.movies)) {
		return []models.Movie{}
	}

	if limit == 0 {
		limit = uint(len(db.movies)) - offset
	}

	result := make(chan []models.Movie)

	go func() {
		defer close(result)
		db.mu.RLock()
		defer db.mu.RUnlock()

		// Simulate fetching movies
		result <- db.movies[offset:min(offset+limit, uint(len(db.movies)))]
	}()
	return <-result
}

func (db *DataBase) FindActorByID(actorID string) *models.Actor {
	result := make(chan *models.Actor)

	go func() {
		defer close(result)
		db.mu.RLock()
		defer db.mu.RUnlock()

		for _, actor := range db.actors {
			if actor.ID == actorID {
				result <- &actor
				return
			}
		}
	}()
	return <-result
}

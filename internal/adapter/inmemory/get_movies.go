package inmemory

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
)

// Get movies from the database
// If limit is 0, return all movies from offset
func (db *DataBase) GetMovies(ctx context.Context, offset, limit uint) ([]entity.Movie, error) {
	db.logger.Info("GetMovies called",
		db.logger.ToString("requestID", ctx.Value(common.RequestIDContextKey).(string)),
		db.logger.ToInt("offset", int(offset)), db.logger.ToInt("limit", int(limit)))

	if offset > uint(len(db.movies)) {
		return []entity.Movie{}, nil
	}

	if limit == 0 {
		limit = uint(len(db.movies)) - offset
	}

	result := make(chan []entity.Movie)

	go func() {
		defer close(result)
		db.mu.RLock()
		defer db.mu.RUnlock()

		// Simulate fetching movies
		result <- db.movies[offset:min(offset+limit, uint(len(db.movies)))]
	}()
	return <-result, nil
}

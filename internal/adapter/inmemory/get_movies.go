package inmemory

import "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"

// Get movies from the database
// If limit is 0, return all movies from offset
func (db *DataBase) GetMovies(offset, limit uint) ([]entity.Movie, error) {
	if offset > uint(len(db.movies)) {
		return []entity.Movie{}, entity.ErrGetMoviesParamsInvalid
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

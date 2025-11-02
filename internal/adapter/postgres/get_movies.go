package postgres

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
)

// Get movies from the database
// If limit is 0, return all movies from offset
func (db *DataBase) GetMovies(ctx context.Context, offset, limit uint) ([]entity.Movie, error) {
	if limit == 0 {
		limit = 100 // Default limit
	}
	var movies []entity.Movie

	// TODO: Remove RANDOM(). Implement proper recommendation algorithm

	// Retrieve random movies
	query := "SELECT id, title, year FROM movies ORDER BY RANDOM() OFFSET $1 LIMIT $2"
	rows, err := db.conn.Query(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie entity.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

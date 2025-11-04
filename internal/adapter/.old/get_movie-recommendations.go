package postgres

import "context"

// // Get movies from the database
// // If limit is 0, return all movies from offset
// func (db *DataBase) GetMovieRecommendations(ctx context.Context, offset, limit uint) ([]entity.Movie, error) {
// 	requestID, ok := ctx.Value(common.RequestIDContextKey).(string)
// 	if !ok {
// 		db.logger.Warn("GetMovieRecommendations: failed to get requestID from context")
// 		requestID = "unknown"
// 	}
// 	db.logger.Info("GetMovieRecommendations called",
// 		db.logger.ToString("requestID", requestID),
// 		db.logger.ToInt("offset", int(offset)), db.logger.ToInt("limit", int(limit)))

// 	// Check total number of movies
// 	movies_len := 0
// 	query := "SELECT COUNT(*) FROM media WHERE type = 'movie'"
// 	err := db.conn.QueryRow(query).Scan(&movies_len)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// If offset is greater than total movies, return empty slice
// 	if offset > uint(movies_len) {
// 		return []entity.Movie{}, nil
// 	}

// 	// If limit is 0, set it to total movies minus offset
// 	if limit == 0 {
// 		limit = uint(movies_len) - offset
// 	}

// 	var movies []entity.Movie

// 	// TODO: Remove RANDOM(). Implement proper recommendation algorithm

// 	return movies, nil
// }

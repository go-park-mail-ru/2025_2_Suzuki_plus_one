package postgres

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
)

// Get total number of media items in the database
func (db *DataBase) GetMediaCount(ctx context.Context, media_type string) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM media WHERE media_type = $1"
	err := db.conn.QueryRow(query, media_type).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetMedia retrieves a media item by its ID from the database
func (db *DataBase) GetMedia(ctx context.Context, media_id uint) (*entity.Media, error) {
	var media entity.Media

	query := `
	SELECT
		media_id,
		media_type,
		title,
		description,
		release_date,
		rating,
		duration_minutes,
		age_rating,
		country,
		plot_summary
	FROM media
	WHERE media_id = $1
	`

	// scan into nullable types to avoid nil/NULL panics
	var mediaID int64
	var mediaType string
	var title string
	var description sql.NullString
	var releaseDate sql.NullTime
	var rating sql.NullFloat64
	var duration sql.NullInt64
	var ageRating sql.NullInt64
	var country sql.NullString
	var plot sql.NullString

	err := db.conn.QueryRow(query, media_id).Scan(
		&mediaID,
		&mediaType,
		&title,
		&description,
		&releaseDate,
		&rating,
		&duration,
		&ageRating,
		&country,
		&plot,
	)
	if err != nil {
		return &entity.Media{}, err
	}

	media.ID = int(mediaID)
	media.MediaType = mediaType
	media.Title = title
	if description.Valid {
		media.Description = description.String
	}
	if releaseDate.Valid {
		media.ReleaseDate = releaseDate.Time
	}
	if rating.Valid {
		media.Rating = rating.Float64
	}
	if duration.Valid {
		media.Duration = int(duration.Int64)
	}
	if ageRating.Valid {
		media.AgeRating = int(ageRating.Int64)
	}
	if country.Valid {
		media.Country = country.String
	}
	if plot.Valid {
		media.PlotSummary = plot.String
	}

	return &media, nil
}

// Get posters s3 links for the given media
func (db *DataBase) GetMediaPostersLinks(ctx context.Context, media_id uint) ([]string, error) {
	var posters []string
	query := `
		SELECT s3_key
		FROM media_image
		JOIN asset_image USING (asset_image_id)
		JOIN asset USING (asset_id)
		WHERE media_id = $1 AND image_type = 'poster'
	`
	rows, err := db.conn.Query(query, media_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		posters = append(posters, url)
	}
	return posters, nil
}

// Get genres for the given media
func (db *DataBase) GetMediaGenres(ctx context.Context, media_id uint) ([]entity.Genre, error) {
	var genres []entity.Genre
	query := `
		SELECT genre_id, name, description
		FROM media_genre
		JOIN genre USING (genre_id)
		WHERE media_id = $1
	`
	rows, err := db.conn.Query(query, media_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var genre entity.Genre
		if err := rows.Scan(&genre.ID, &genre.Name, &genre.Description); err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	return genres, nil
}

package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
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

// GetMediaByID retrieves a media item by its ID from the database
func (db *DataBase) GetMediaByID(ctx context.Context, media_id uint) (*entity.Media, error) {
	// Log the request ID from context for tracing
	requestID, ok := ctx.Value(common.RequestIDContextKey).(string)
	if !ok {
		db.logger.Warn("GetMediaByID: failed to get requestID from context")
		requestID = "unknown"
	}
	db.logger.Info("GetMediaByID called",
		db.logger.ToString("requestID", requestID),
		db.logger.ToInt("media_id", int(media_id)),
	)

	// Prepare the media entity to be filled
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

	media.MediaID = uint(mediaID)
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

// Get posters s3 keys for the given media
func (db *DataBase) GetMediaPostersKeys(ctx context.Context, media_id uint) ([]entity.S3Key, error) {
	var posters []entity.S3Key
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
		// Split bucket name and key from url
		s3key, err := splitS3Key(url)
		if err == nil {
			posters = append(posters, s3key)
		} else {
			db.logger.Error("GetMediaPostersKeys: failed to split S3 key", db.logger.ToError(err))
		}
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

// Get random media IDs for recommendations using RANDOM()
func (db *DataBase) GetMediaRandomIds(ctx context.Context, limit uint, offset uint, media_type string) ([]uint, error) {
	var mediaIDs []uint
	query := `
		SELECT media_id
		FROM media
		WHERE media_type = $1
		ORDER BY RANDOM()
		LIMIT $2 OFFSET $3
	`
	rows, err := db.conn.Query(query, media_type, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var mediaID uint
		if err := rows.Scan(&mediaID); err != nil {
			return nil, err
		}
		mediaIDs = append(mediaIDs, mediaID)
	}
	return mediaIDs, nil
}

func (db *DataBase) GetMediaWatchKey(ctx context.Context, media_id uint) (*entity.S3Key, error) {
	var s3Key string
	query := `
		SELECT s3_key
		FROM media
		JOIN media_video USING (media_id)
		JOIN asset_video USING (asset_video_id)
		JOIN asset USING (asset_id)
		WHERE media_id = $1
	`
	err := db.conn.QueryRow(query, media_id).Scan(&s3Key)
	if err != nil {
		return nil, err
	}
	// Split bucket name and key from s3Key
	s3key, err := splitS3Key(s3Key)
	if err != nil {
		db.logger.Error("GetMediaWatchKey: failed to split S3 key", db.logger.ToError(err))
		return nil, err
	}
	return &s3key, nil
}

var ErrInvalidS3KeyFormat = errors.New("invalid S3 key format")

// Splits ""bucket_name/key"" into S3Key struct
func splitS3Key(fullPath string) (entity.S3Key, error) {
	// fullPath cannot start with '/'
	if fullPath[0] == '/' {
		fullPath = fullPath[1:]
	}

	var bucketName, key string
	splitIndex := -1
	for i, char := range fullPath {
		if char == '/' {
			splitIndex = i
			break
		}
	}
	if splitIndex != -1 {
		bucketName = fullPath[:splitIndex]
		key = fullPath[splitIndex+1:]
	} else {
		return entity.S3Key{}, ErrInvalidS3KeyFormat
	}
	return entity.S3Key{
		BucketName: bucketName,
		Key:        key,
	}, nil
}

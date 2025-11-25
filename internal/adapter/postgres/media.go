package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// Get total number of media items in the database
func (db *DataBase) GetMediaCount(ctx context.Context, media_type string) (int, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetMediaCount called",
		log.ToString("media_type", media_type),
	)
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
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetMediaByID called",
		log.ToInt("media_id", int(media_id)),
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
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetMediaPostersKeys called",
		log.ToInt("media_id", int(media_id)),
	)

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
			log.Error("GetMediaPostersKeys: failed to scan URL", log.ToError(err))
			return nil, err
		}
		// Split bucket name and key from url
		s3key, err := splitS3Key(url)
		if err == nil {
			posters = append(posters, s3key)
		} else {
			log.Error("GetMediaPostersKeys: failed to split S3 key", log.ToError(err))
		}
	}
	return posters, nil
}

// Get genres for the given media
func (db *DataBase) GetMediaGenres(ctx context.Context, media_id uint) ([]entity.Genre, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetMediaGenres called",
		log.ToInt("media_id", int(media_id)),
	)

	var genres []entity.Genre
	query := `
		SELECT genre_id, name, description
		FROM media_genre
		JOIN genre USING (genre_id)
		WHERE media_id = $1
	`
	rows, err := db.conn.Query(query, media_id)
	if err != nil {
		log.Error("GetMediaGenres: failed to execute query", log.ToError(err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var genre entity.Genre
		if err := rows.Scan(&genre.ID, &genre.Name, &genre.Description); err != nil {
			log.Error("GetMediaGenres: failed to scan genre", log.ToError(err))
			return nil, err
		}
		genres = append(genres, genre)
	}
	return genres, nil
}

// Get media trailer keys for the given media
func (db *DataBase) GetMediaTrailersKeys(ctx context.Context, media_id uint) ([]entity.S3Key, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetMediaTrailerKeys called",
		log.ToInt("media_id", int(media_id)),
	)

	var trailerKeys []entity.S3Key
	query := `
		SELECT s3_key
		FROM MEDIA_VIDEO
		JOIN ASSET_VIDEO USING (asset_video_id)
		JOIN ASSET USING (asset_id)
		WHERE media_id = $1 AND video_type = 'trailer'
	`
	rows, err := db.conn.Query(query, media_id)
	if err != nil {
		log.Error("GetMediaTrailerKeys: failed to execute query", log.ToError(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s3KeyString string
		if err := rows.Scan(&s3KeyString); err != nil {
			log.Error("GetMediaTrailersKeys: failed to scan S3 key", log.ToError(err))
			return nil, err
		}
		s3Key, err := splitS3Key(s3KeyString)
		if err != nil {
			log.Error("GetMediaTrailersKeys: failed to split S3 key", log.ToError(err))
			return nil, err
		}
		trailerKeys = append(trailerKeys, s3Key)
	}
	return trailerKeys, nil
}

// Get random media IDs for recommendations using RANDOM()
func (db *DataBase) GetMediaSortedByName(ctx context.Context, limit uint, offset uint, media_type string) ([]uint, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetMediaSortedByName called",
		log.ToInt("limit", int(limit)),
		log.ToInt("offset", int(offset)),
		log.ToString("media_type", media_type),
	)
	var mediaIDs []uint
	query := `
		SELECT media_id
		FROM media
		WHERE media_type = $1
		ORDER BY title
		LIMIT $2 OFFSET $3
	`
	rows, err := db.conn.Query(query, media_type, limit, offset)
	if err != nil {
		log.Error("GetMediaSortedByName: failed to execute query", log.ToError(err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var mediaID uint
		if err := rows.Scan(&mediaID); err != nil {
			log.Error("GetMediaSortedByName: failed to scan media ID", log.ToError(err))
			return nil, err
		}
		mediaIDs = append(mediaIDs, mediaID)
	}
	return mediaIDs, nil
}

func (db *DataBase) GetMediaWatchKey(ctx context.Context, media_id uint) (*entity.S3Key, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetMediaWatchKey called",
		log.ToInt("media_id", int(media_id)),
	)

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
		log.Error("GetMediaWatchKey: failed to execute query", log.ToError(err))
		return nil, err
	}

	// Split bucket name and key from s3Key
	s3key, err := splitS3Key(s3Key)
	if err != nil {
		log.Error("GetMediaWatchKey: failed to split S3 key", log.ToError(err))
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


func (db *DataBase) GetMediaIDsByLikeStatus(ctx context.Context, userID uint, isDislike bool, limit uint, offset uint) ([]uint, error) {
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetMediaIDsByLikeStatus called",
		log.ToInt("user_id", int(userID)),
		log.ToAny("is_dislike", isDislike),
		log.ToInt("limit", int(limit)),
		log.ToInt("offset", int(offset)),
	)

	var mediaIDs []uint
	query := `
	SELECT media_id
	FROM user_like_media
	WHERE user_id = $1 AND is_dislike = $2
	LIMIT $3 OFFSET $4
	`
	rows, err := db.conn.Query(query, userID, isDislike, limit, offset)
	if err != nil {
		log.Error("GetMediaIDsByLikeStatus: failed to execute query", log.ToError(err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var mediaID uint
		if err := rows.Scan(&mediaID); err != nil {
			log.Error("GetMediaIDsByLikeStatus: failed to scan media ID", log.ToError(err))
			return nil, err
		}
		mediaIDs = append(mediaIDs, mediaID)
	}
	return mediaIDs, nil
}
package postgres

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// Get actor by their ID
func (db *DataBase) GetActorByID(ctx context.Context, actorID uint) (*entity.Actor, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContexKeyRequestID)

	var actor entity.Actor
	query := "SELECT actor_id, name, bio, birth_date FROM actor WHERE actor_id = $1"
	err := db.conn.QueryRowContext(ctx, query, actorID).Scan(
		&actor.ID,
		&actor.Name,
		&actor.Bio,
		&actor.BirthDate,
	)
	if err != nil {
		log.Error("Failed to get actor by ID: " + err.Error())
		return nil, err
	}
	return &actor, nil
}

// Get actor images S3 keys
func (db *DataBase) GetActorImageS3(ctx context.Context, actorID uint) ([]entity.S3Key, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContexKeyRequestID)

	var imageURLs []entity.S3Key
	query := `
		SELECT s3_key
		FROM actor
		JOIN actor_image USING (actor_id)
		JOIN asset_image USING (asset_image_id)
		JOIN asset USING (asset_id)
		WHERE actor_id = $1
	`
	rows, err := db.conn.QueryContext(ctx, query, actorID)
	if err != nil {
		log.Error("Failed to query actor images: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			log.Error("Failed to scan actor image URL: " + err.Error())
			return nil, err
		}
		s3Key, err := splitS3Key(url)
		if err == nil {
			imageURLs = append(imageURLs, s3Key)
		} else {
			log.Error("GetActorImageS3: failed to split S3 key", log.ToError(err))
		}
	}
	return imageURLs, nil
}

// Get medias where actor had a role
func (db *DataBase) GetMediasByActorID(ctx context.Context, actorID uint) ([]entity.Media, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContexKeyRequestID)

	log.Info("GetMediasByActorID called",
		log.ToInt("actor_id", int(actorID)),
	)

	var medias []entity.Media
	query := `
		SELECT media.media_id, media_type, title, description
		FROM media
		JOIN actor_role USING (media_id)
		WHERE actor_id = $1
	`
	rows, err := db.conn.QueryContext(ctx, query, actorID)
	if err != nil {
		log.Error("Failed to query medias by actor ID: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var media entity.Media
		var description sql.NullString
		if err := rows.Scan(&media.MediaID, &media.MediaType, &media.Title, &description); err != nil {
			log.Error("Failed to scan media: " + err.Error())
			return nil, err
		}
		if description.Valid {
			media.Description = description.String
		}
		medias = append(medias, media)
	}
	return medias, nil
}

// Get actors for the given media
func (db *DataBase) GetActorsByMediaID(ctx context.Context, media_id uint) ([]entity.Actor, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContexKeyRequestID)

	log.Info("GetActorsByMediaID called",
		log.ToInt("media_id", int(media_id)),
	)

	var actors []entity.Actor
	query := `
		SELECT actor_id, name, birth_date, bio
		FROM actor
		JOIN actor_role USING (actor_id)
		WHERE media_id = $1
	`
	rows, err := db.conn.Query(query, media_id)
	if err != nil {
		log.Error("Failed to query actors by media ID: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var actor entity.Actor
		if err := rows.Scan(&actor.ID, &actor.Name, &actor.BirthDate, &actor.Bio); err != nil {
			log.Error("Failed to scan actor: " + err.Error())
			return nil, err
		}
		actors = append(actors, actor)
	}
	return actors, nil
}

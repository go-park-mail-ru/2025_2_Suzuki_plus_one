package postgres

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
)

// Get actor by their ID
func (db *DataBase) GetActorByID(ctx context.Context, actorID uint) (*entity.Actor, error) {
	var actor entity.Actor
	query := "SELECT actor_id, name, bio, birth_date FROM actor WHERE actor_id = $1"
	err := db.conn.QueryRowContext(ctx, query, actorID).Scan(
		&actor.ID,
		&actor.Name,
		&actor.Bio,
		&actor.BirthDate,
	)
	if err != nil {
		return nil, err
	}
	return &actor, nil
}

// Get actor images S3 keys
func (db *DataBase) GetActorImageS3(ctx context.Context, actorID uint) ([]string, error) {
	var imageURLs []string
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
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		imageURLs = append(imageURLs, url)
	}
	return imageURLs, nil
}

// Get medias where actor had a role
func (db *DataBase) GetMediasByActorID(ctx context.Context, actorID uint) ([]entity.Media, error) {
	var medias []entity.Media
	query := `
		SELECT media.media_id, media_type, title, description
		FROM media
		JOIN actor_role USING (media_id)
		WHERE actor_id = $1
	`
	rows, err := db.conn.QueryContext(ctx, query, actorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var media entity.Media
		var description sql.NullString
		if err := rows.Scan(&media.MediaID, &media.MediaType, &media.Title, &description); err != nil {
			return nil, err
		}
		if description.Valid {
			media.Description = description.String
		}
		medias = append(medias, media)
	}
	return medias, nil
}

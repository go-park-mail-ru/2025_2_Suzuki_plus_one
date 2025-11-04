package postgres

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
)

func (db *DataBase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User

	query := `
		SELECT user_id, email, username, password_hash, asset_image_id
		FROM "user"
		WHERE email = $1
	`
	row := db.conn.QueryRow(query, email)
	var assetImageID sql.NullInt64

	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &assetImageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
	}
	if assetImageID.Valid {
		user.AssetImageID = uint(assetImageID.Int64)
	}

	return &user, nil
}

func (db *DataBase) GetUserByID(ctx context.Context, userID uint) (*entity.User, error) {
	var user entity.User

	query := `
		SELECT user_id, email, username, password_hash, asset_image_id
		FROM "user"
		WHERE user_id = $1
	`
	row := db.conn.QueryRow(query, userID)
	var assetImageID sql.NullInt64
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &assetImageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
	}
	if assetImageID.Valid {
		user.AssetImageID = uint(assetImageID.Int64)
	}

	return &user, nil
}

// Returns the S3 key (URL) of the user's avatar image
func (db *DataBase) GetUserAvatarURL(ctx context.Context, userID uint) (string, error) {
	var avatarURL string

	query := `
		SELECT s3_key
		FROM "user"
		JOIN asset_image USING (asset_image_id)
		JOIN asset USING (asset_id)
		WHERE user_id = $1
	`
	row := db.conn.QueryRow(query, userID)
	err := row.Scan(&avatarURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", entity.ErrUserNotFound
		}
		return "", err
	}

	return avatarURL, nil
}

func (db *DataBase) CreateUser(ctx context.Context, user entity.User) (uint, error) {
	var userID uint

	query := `
		INSERT INTO "user" (email, username, password_hash)
		VALUES ($1, $2, $3)
		RETURNING user_id
	`
	err := db.conn.QueryRow(query, user.Email, user.Username, user.PasswordHash).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

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
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.AssetImageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
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
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.AssetImageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
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

package postgres

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

func (db *DataBase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetUserByEmail called",
		log.ToString("email", email),
	)

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
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetUserByID called",
		log.ToInt("user_id", int(userID)),
	)

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

// Returns the S3 key of the user's avatar image
func (db *DataBase) GetUserAvatarKey(ctx context.Context, userID uint) (*entity.S3Key, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetUserAvatarKey called",
		log.ToInt("user_id", int(userID)),
	)

	var avatarKey string

	query := `
		SELECT s3_key
		FROM "user"
		JOIN asset_image USING (asset_image_id)
		JOIN asset USING (asset_id)
		WHERE user_id = $1
	`
	row := db.conn.QueryRow(query, userID)
	err := row.Scan(&avatarKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
	}

	s3Key, err := splitS3Key(avatarKey)
	if err != nil {
		db.logger.Error("GetUserAvatarKey: failed to split S3 key", db.logger.ToError(err))
		return nil, err
	}
	return &s3Key, nil
}

func (db *DataBase) CreateUser(ctx context.Context, user entity.User) (uint, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("CreateUser called",
		log.ToString("email", user.Email),
		log.ToString("username", user.Username),
	)

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

func (db *DataBase) UpdateUser(
	ctx context.Context,
	userID uint,
	username string,
	email string,
	dateOfBirth string,
	phoneNumber string,
) (*entity.User, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("UpdateUser called",
		log.ToInt("user_id", int(userID)),
		log.ToString("username", username),
		log.ToString("email", email),
		log.ToString("date_of_birth", dateOfBirth),
		log.ToString("phone_number", phoneNumber),
	)

	query := `
		UPDATE "user"
		SET username = $1,
		    email = $2,
		    date_of_birth = $3,
		    phone_number = $4
		WHERE user_id = $5
		RETURNING user_id, email, username, date_of_birth, phone_number, password_hash, asset_image_id
	`
	row := db.conn.QueryRow(query, username, email, dateOfBirth, phoneNumber, userID)

	var updatedUser entity.User
	var assetImageID sql.NullInt64
	err := row.Scan(&updatedUser.ID, &updatedUser.Email, &updatedUser.Username, &updatedUser.DateOfBirth, &updatedUser.PhoneNumber, &updatedUser.PasswordHash, &assetImageID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error("UpdateUser: user not found", log.ToInt("user_id", int(userID)))
			return nil, entity.ErrUserNotFound
		}
		return nil, err
	}
	if assetImageID.Valid {
		updatedUser.AssetImageID = uint(assetImageID.Int64)
	}

	return &updatedUser, nil
}

func (db *DataBase) UpdateUserAvatarKey(ctx context.Context, userID uint, assetImageID uint) error {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("UpdateUserAvatarKey called",
		log.ToInt("user_id", int(userID)),
		log.ToInt("asset_image_id", int(assetImageID)),
	)

	query := `
		UPDATE "user"
		SET asset_image_id = $1
		WHERE user_id = $2
	`
	result, err := db.conn.Exec(query, assetImageID, userID)
	if err != nil {
		log.Error("UpdateUserAvatarKey: failed to update user avatar key", log.ToError(err))
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("UpdateUserAvatarKey: failed to get rows affected", log.ToError(err))
		return err
	}
	if rowsAffected == 0 {
		log.Error("UpdateUserAvatarKey: user not found", log.ToInt("user_id", int(userID)))
		return entity.ErrUserNotFound
	}

	return nil
}

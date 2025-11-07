package postgres

import (
	"context"
	"database/sql"
	"time"

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
		SELECT user_id, email, username, password_hash, asset_image_id, date_of_birth, phone_number
		FROM "user"
		WHERE email = $1
	`
	row := db.conn.QueryRow(query, email)
	err := row.Scan(&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.AssetImageID,
		&user.DateOfBirth,
		&user.PhoneNumber,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error("GetUserByEmail: user not found", log.ToString("email", email))
			return nil, entity.ErrUserNotFound
		}
		log.Error("GetUserByEmail: failed to scan user", log.ToError(err))
		return nil, err
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
		SELECT user_id, email, username, password_hash, asset_image_id, date_of_birth, phone_number
		FROM "user"
		WHERE user_id = $1
	`
	row := db.conn.QueryRow(query, userID)
	err := row.Scan(&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.AssetImageID,
		&user.DateOfBirth,
		&user.PhoneNumber,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error("GetUserByID: user not found", log.ToInt("user_id", int(userID)))
			return nil, entity.ErrUserNotFound
		}
		log.Error("GetUserByID: failed to scan user", log.ToError(err))
		return nil, err
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

func (db *DataBase) CreateUser(ctx context.Context, email string, username string, passwordHash string) (uint, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("CreateUser called",
		log.ToString("email", email),
		log.ToString("username", username),
	)

	var userID uint

	query := `
		INSERT INTO "user" (email, username, password_hash)
		VALUES ($1, $2, $3)
		RETURNING user_id
	`
	err := db.conn.QueryRow(query,
		email,
		username,
		passwordHash,
	).Scan(&userID)
	if err != nil {
		log.Error("CreateUser: failed to create user", log.ToError(err))
		return 0, err
	}

	return userID, nil
}

func (db *DataBase) UpdateUser(
	ctx context.Context,
	userID uint,
	username string,
	email string,
	dateOfBirth time.Time,
	phoneNumber string,
) (*entity.User, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("UpdateUser called",
		log.ToInt("user_id", int(userID)),
		log.ToString("username", username),
		log.ToString("email", email),
		log.ToAny("date_of_birth", dateOfBirth),
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

	err := row.Scan(
		&updatedUser.ID,
		&updatedUser.Email,
		&updatedUser.Username,
		&updatedUser.DateOfBirth,
		&updatedUser.PhoneNumber,
		&updatedUser.PasswordHash,
		&updatedUser.AssetImageID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error("UpdateUser: user not found", log.ToInt("user_id", int(userID)))
			return nil, entity.ErrUserNotFound
		}
		return nil, err
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

func (db *DataBase) UpdateUserPassword(ctx context.Context, userID uint, newHashedPassword string) error {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("UpdateUserPassword called",
		log.ToInt("user_id", int(userID)),
	)

	query := `
		UPDATE "user"
		SET password_hash = $1
		WHERE user_id = $2
	`
	result, err := db.conn.Exec(query, newHashedPassword, userID)
	if err != nil {
		log.Error("UpdateUserPassword: failed to update user password", log.ToError(err))
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("UpdateUserPassword: failed to get rows affected", log.ToError(err))
		return err
	}
	if rowsAffected == 0 {
		log.Error("UpdateUserPassword: user not found", log.ToInt("user_id", int(userID)))
		return entity.ErrUserNotFound
	}

	return nil
}

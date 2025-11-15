package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	logpkg "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserByEmail_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.RequestIDContextKey, "test-request-id")

	db, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to open sqlmock database")
	defer db.Close()

	lg := logpkg.NewZapLogger(true)

	database := &DataBase{
		conn:    db,
		logger:  lg,
		context: ctx,
	}

	email := "test@example.com"
	expectedUser := &entity.User{
		ID:           123,
		Email:        email,
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		AssetImageID: 456,
	}

	rows := sqlmock.NewRows([]string{"user_id", "email", "username", "password_hash", "asset_image_id"}).
		AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Username, expectedUser.PasswordHash, expectedUser.AssetImageID)

	mock.ExpectQuery(`SELECT user_id, email, username, password_hash, asset_image_id FROM "user" WHERE email = \$1`).
		WithArgs(email).
		WillReturnRows(rows)

	user, err := database.GetUserByEmail(ctx, email)
	require.NoError(t, err)
	require.NotNil(t, user)

	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.PasswordHash, user.PasswordHash)
	assert.Equal(t, expectedUser.AssetImageID, user.AssetImageID)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetUserByEmail_WithNullAssetImageID(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.RequestIDContextKey, "test-request-id")

	db, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to open sqlmock database")
	defer db.Close()

	lg := logpkg.NewZapLogger(true)

	database := &DataBase{
		conn:    db,
		logger:  lg,
		context: ctx,
	}

	email := "test@example.com"

	rows := sqlmock.NewRows([]string{"user_id", "email", "username", "password_hash", "asset_image_id"}).
		AddRow(123, email, "testuser", "hashedpassword", nil)

	mock.ExpectQuery(`SELECT user_id, email, username, password_hash, asset_image_id FROM "user" WHERE email = \$1`).
		WithArgs(email).
		WillReturnRows(rows)

	user, err := database.GetUserByEmail(ctx, email)
	require.NoError(t, err)
	require.NotNil(t, user)

	assert.Equal(t, uint(123), user.ID)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "hashedpassword", user.PasswordHash)
	assert.Equal(t, uint(0), user.AssetImageID)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.RequestIDContextKey, "test-request-id")

	db, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to open sqlmock database")
	defer db.Close()

	lg := logpkg.NewZapLogger(true)

	database := &DataBase{
		conn:    db,
		logger:  lg,
		context: ctx,
	}

	email := "nonexistent@example.com"

	mock.ExpectQuery(`SELECT user_id, email, username, password_hash, asset_image_id FROM "user" WHERE email = \$1`).
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	user, err := database.GetUserByEmail(ctx, email)
	require.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, entity.ErrUserNotFound, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetUserByEmail_DBError(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.RequestIDContextKey, "test-request-id")

	db, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to open sqlmock database")
	defer db.Close()

	lg := logpkg.NewZapLogger(true)

	database := &DataBase{
		conn:    db,
		logger:  lg,
		context: ctx,
	}

	email := "test@example.com"

	mock.ExpectQuery(`SELECT user_id, email, username, password_hash, asset_image_id FROM "user" WHERE email = \$1`).
		WithArgs(email).
		WillReturnError(errors.New("database error"))

	user, err := database.GetUserByEmail(ctx, email)
	require.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "database error")

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetUserByID_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.RequestIDContextKey, "test-request-id")

	db, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to open sqlmock database")
	defer db.Close()

	lg := logpkg.NewZapLogger(true)

	database := &DataBase{
		conn:    db,
		logger:  lg,
		context: ctx,
	}

	userID := uint(123)
	expectedUser := &entity.User{
		ID:           userID,
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		AssetImageID: 456,
	}

	rows := sqlmock.NewRows([]string{"user_id", "email", "username", "password_hash", "asset_image_id"}).
		AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Username, expectedUser.PasswordHash, expectedUser.AssetImageID)

	mock.ExpectQuery(`SELECT user_id, email, username, password_hash, asset_image_id FROM "user" WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(rows)

	user, err := database.GetUserByID(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, user)

	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.PasswordHash, user.PasswordHash)
	assert.Equal(t, expectedUser.AssetImageID, user.AssetImageID)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetUserByID_NotFound(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.RequestIDContextKey, "test-request-id")

	db, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to open sqlmock database")
	defer db.Close()

	lg := logpkg.NewZapLogger(true)

	database := &DataBase{
		conn:    db,
		logger:  lg,
		context: ctx,
	}

	userID := uint(999)

	mock.ExpectQuery(`SELECT user_id, email, username, password_hash, asset_image_id FROM "user" WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	user, err := database.GetUserByID(ctx, userID)
	require.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, entity.ErrUserNotFound, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetUserAvatarKey_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.RequestIDContextKey, "test-request-id")

	db, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to open sqlmock database")
	defer db.Close()

	lg := logpkg.NewZapLogger(true)

	database := &DataBase{
		conn:    db,
		logger:  lg,
		context: ctx,
	}

	userID := uint(123)
	expectedAvatarKey := "avatars/user123.jpg"

	rows := sqlmock.NewRows([]string{"s3_key"}).
		AddRow(expectedAvatarKey)

	mock.ExpectQuery(`SELECT s3_key FROM "user" JOIN asset_image USING \(asset_image_id\) JOIN asset USING \(asset_id\) WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(rows)

	avatarKey, err := database.GetUserAvatarKey(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, expectedAvatarKey, avatarKey)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetUserAvatarKey_NotFound(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.RequestIDContextKey, "test-request-id")

	db, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to open sqlmock database")
	defer db.Close()

	lg := logpkg.NewZapLogger(true)

	database := &DataBase{
		conn:    db,
		logger:  lg,
		context: ctx,
	}

	userID := uint(999)

	mock.ExpectQuery(`SELECT s3_key FROM "user" JOIN asset_image USING \(asset_image_id\) JOIN asset USING \(asset_id\) WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	avatarKey, err := database.GetUserAvatarKey(ctx, userID)
	require.Error(t, err)
	assert.Equal(t, "", avatarKey)
	assert.Equal(t, entity.ErrUserNotFound, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetUserAvatarKey_DBError(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.RequestIDContextKey, "test-request-id")

	db, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to open sqlmock database")
	defer db.Close()

	lg := logpkg.NewZapLogger(true)

	database := &DataBase{
		conn:    db,
		logger:  lg,
		context: ctx,
	}

	userID := uint(123)

	mock.ExpectQuery(`SELECT s3_key FROM "user" JOIN asset_image USING \(asset_image_id\) JOIN asset USING \(asset_id\) WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnError(errors.New("join failed"))

	avatarKey, err := database.GetUserAvatarKey(ctx, userID)
	require.Error(t, err)
	assert.Equal(t, "", avatarKey)
	assert.Contains(t, err.Error(), "join failed")

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestCreateUser_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.RequestIDContextKey, "test-request-id")

	db, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to open sqlmock database")
	defer db.Close()

	lg := logpkg.NewZapLogger(true)

	database := &DataBase{
		conn:    db,
		logger:  lg,
		context: ctx,
	}

	user := entity.User{
		Email:        "newuser@example.com",
		Username:     "newuser",
		PasswordHash: "newhashedpassword",
	}
	expectedUserID := uint(456)

	rows := sqlmock.NewRows([]string{"user_id"}).AddRow(expectedUserID)

	mock.ExpectQuery(`INSERT INTO "user" \(email, username, password_hash\) VALUES \(\$1, \$2, \$3\) RETURNING user_id`).
		WithArgs(user.Email, user.Username, user.PasswordHash).
		WillReturnRows(rows)

	userID, err := database.CreateUser(ctx, user)
	require.NoError(t, err)
	assert.Equal(t, expectedUserID, userID)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.RequestIDContextKey, "test-request-id")

	db, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to open sqlmock database")
	defer db.Close()

	lg := logpkg.NewZapLogger(true)

	database := &DataBase{
		conn:    db,
		logger:  lg,
		context: ctx,
	}

	user := entity.User{
		Email:        "existing@example.com",
		Username:     "newuser",
		PasswordHash: "newhashedpassword",
	}

	mock.ExpectQuery(`INSERT INTO "user" \(email, username, password_hash\) VALUES \(\$1, \$2, \$3\) RETURNING user_id`).
		WithArgs(user.Email, user.Username, user.PasswordHash).
		WillReturnError(errors.New("duplicate key value violates unique constraint"))

	userID, err := database.CreateUser(ctx, user)
	require.Error(t, err)
	assert.Equal(t, uint(0), userID)
	assert.Contains(t, err.Error(), "duplicate")

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestCreateUser_DBError(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.RequestIDContextKey, "test-request-id")

	db, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to open sqlmock database")
	defer db.Close()

	lg := logpkg.NewZapLogger(true)

	database := &DataBase{
		conn:    db,
		logger:  lg,
		context: ctx,
	}

	user := entity.User{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}

	mock.ExpectQuery(`INSERT INTO "user" \(email, username, password_hash\) VALUES \(\$1, \$2, \$3\) RETURNING user_id`).
		WithArgs(user.Email, user.Username, user.PasswordHash).
		WillReturnError(errors.New("connection failed"))

	userID, err := database.CreateUser(ctx, user)
	require.Error(t, err)
	assert.Equal(t, uint(0), userID)
	assert.Contains(t, err.Error(), "connection failed")

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestUserMethods_EdgeCases(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.RequestIDContextKey, "test-request-id")

	db, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to open sqlmock database")
	defer db.Close()

	lg := logpkg.NewZapLogger(true)

	database := &DataBase{
		conn:    db,
		logger:  lg,
		context: ctx,
	}

	t.Run("ZeroUserID", func(t *testing.T) {
		mock.ExpectQuery(`SELECT user_id, email, username, password_hash, asset_image_id FROM "user" WHERE user_id = \$1`).
			WithArgs(uint(0)).
			WillReturnError(sql.ErrNoRows)

		user, err := database.GetUserByID(ctx, 0)
		require.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, entity.ErrUserNotFound, err)
	})

	t.Run("EmptyEmail", func(t *testing.T) {
		mock.ExpectQuery(`SELECT user_id, email, username, password_hash, asset_image_id FROM "user" WHERE email = \$1`).
			WithArgs("").
			WillReturnError(sql.ErrNoRows)

		user, err := database.GetUserByEmail(ctx, "")
		require.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, entity.ErrUserNotFound, err)
	})

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

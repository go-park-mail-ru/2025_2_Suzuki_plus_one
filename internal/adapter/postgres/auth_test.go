package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	logpkg "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddNewRefreshToken_Success(t *testing.T) {
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
	refreshToken := "test-refresh-token"
	expiresAt := time.Now().Add(24 * time.Hour)

	mock.ExpectExec("INSERT INTO user_session").
		WithArgs(userID, refreshToken, expiresAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = database.AddNewRefreshToken(ctx, userID, refreshToken, expiresAt)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestAddNewRefreshToken_NoRequestID(t *testing.T) {
	ctx := context.Background()

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
	refreshToken := "test-refresh-token"
	expiresAt := time.Now().Add(24 * time.Hour)

	mock.ExpectExec("INSERT INTO user_session").
		WithArgs(userID, refreshToken, expiresAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = database.AddNewRefreshToken(ctx, userID, refreshToken, expiresAt)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestAddNewRefreshToken_DBError(t *testing.T) {
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
	refreshToken := "test-refresh-token"
	expiresAt := time.Now().Add(24 * time.Hour)
	
	mock.ExpectExec("INSERT INTO user_session").
		WithArgs(userID, refreshToken, expiresAt).
		WillReturnError(errors.New("database connection failed"))

	err = database.AddNewRefreshToken(ctx, userID, refreshToken, expiresAt)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "database connection failed")

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetRefreshTokensForUser_Success(t *testing.T) {
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
	expiresAt := time.Now().Add(24 * time.Hour)

	rows := sqlmock.NewRows([]string{"user_session_id", "session_token", "expires_at"}).
		AddRow(1, "token-1", expiresAt).
		AddRow(2, "token-2", expiresAt.Add(time.Hour)).
		AddRow(3, "token-3", expiresAt.Add(2*time.Hour))

	mock.ExpectQuery(`SELECT user_session_id, session_token, expires_at FROM user_session WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(rows)

	tokens, err := database.GetRefreshTokensForUser(ctx, userID)
	require.NoError(t, err)
	require.Len(t, tokens, 3)

	assert.Equal(t, uint(1), tokens[0].ID)
	assert.Equal(t, "token-1", tokens[0].Token)
	assert.Equal(t, uint(0), tokens[0].UserID)
	assert.Equal(t, expiresAt, tokens[0].ExpiresAt)

	assert.Equal(t, uint(2), tokens[1].ID)
	assert.Equal(t, "token-2", tokens[1].Token)
	assert.Equal(t, uint(0), tokens[1].UserID)
	assert.Equal(t, expiresAt.Add(time.Hour), tokens[1].ExpiresAt)

	assert.Equal(t, uint(3), tokens[2].ID)
	assert.Equal(t, "token-3", tokens[2].Token)
	assert.Equal(t, uint(0), tokens[2].UserID)
	assert.Equal(t, expiresAt.Add(2*time.Hour), tokens[2].ExpiresAt)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetRefreshTokensForUser_NoTokens(t *testing.T) {
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

	rows := sqlmock.NewRows([]string{"user_session_id", "session_token", "expires_at"})

	mock.ExpectQuery("SELECT user_session_id, session_token, expires_at FROM user_session").
		WithArgs(userID).
		WillReturnRows(rows)

	tokens, err := database.GetRefreshTokensForUser(ctx, userID)
	require.NoError(t, err)
	require.Empty(t, tokens)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetRefreshTokensForUser_DBError(t *testing.T) {
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

	mock.ExpectQuery("SELECT user_session_id, session_token, expires_at FROM user_session").
		WithArgs(userID).
		WillReturnError(errors.New("query failed"))

	tokens, err := database.GetRefreshTokensForUser(ctx, userID)
	require.Error(t, err)
	assert.Nil(t, tokens)
	assert.Contains(t, err.Error(), "query failed")

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetRefreshTokensForUser_ScanError(t *testing.T) {
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

	rows := sqlmock.NewRows([]string{"user_session_id", "session_token", "expires_at"}).
		AddRow(1, "token-1", "invalid-date") // Invalid date format

	mock.ExpectQuery("SELECT user_session_id, session_token, expires_at FROM user_session").
		WithArgs(userID).
		WillReturnRows(rows)

	tokens, err := database.GetRefreshTokensForUser(ctx, userID)
	require.Error(t, err)
	assert.Nil(t, tokens)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestRemoveRefreshToken_Success(t *testing.T) {
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
	refreshToken := "token-to-delete"

	mock.ExpectExec("DELETE FROM user_session").
		WithArgs(userID, refreshToken).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = database.RemoveRefreshToken(ctx, userID, refreshToken)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestRemoveRefreshToken_NoRequestID(t *testing.T) {
	ctx := context.Background()

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
	refreshToken := "token-to-delete"

	mock.ExpectExec("DELETE FROM user_session").
		WithArgs(userID, refreshToken).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = database.RemoveRefreshToken(ctx, userID, refreshToken)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestRemoveRefreshToken_DBError(t *testing.T) {
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
	refreshToken := "token-to-delete"

	mock.ExpectExec("DELETE FROM user_session").
		WithArgs(userID, refreshToken).
		WillReturnError(errors.New("delete failed"))

	err = database.RemoveRefreshToken(ctx, userID, refreshToken)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "delete failed")

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestRemoveRefreshToken_NoRowsAffected(t *testing.T) {
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
	refreshToken := "non-existent-token"

	mock.ExpectExec("DELETE FROM user_session").
		WithArgs(userID, refreshToken).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = database.RemoveRefreshToken(ctx, userID, refreshToken)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestAuthMethods_ConcurrentAccess(t *testing.T) {
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
	refreshToken := "test-token"
	expiresAt := time.Now().Add(24 * time.Hour)

	t.Run("AddToken", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO user_session`).
			WithArgs(userID, refreshToken, expiresAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := database.AddNewRefreshToken(ctx, userID, refreshToken, expiresAt)
		assert.NoError(t, err)
	})

	t.Run("GetTokens", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"user_session_id", "session_token", "expires_at"}).
			AddRow(1, refreshToken, expiresAt)

		mock.ExpectQuery(`SELECT user_session_id, session_token, expires_at FROM user_session`).
			WithArgs(userID).
			WillReturnRows(rows)

		tokens, err := database.GetRefreshTokensForUser(ctx, userID)
		assert.NoError(t, err)
		assert.Len(t, tokens, 1)
		assert.Equal(t, refreshToken, tokens[0].Token)
	})

	t.Run("RemoveToken", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM user_session`).
			WithArgs(userID, refreshToken).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := database.RemoveRefreshToken(ctx, userID, refreshToken)
		assert.NoError(t, err)
	})

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

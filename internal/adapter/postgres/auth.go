package postgres

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
)

func (db *DataBase) AddNewRefreshToken(ctx context.Context, userID uint, refreshToken string, expiresAt time.Time) error {
	// Log the request ID from context for tracing
	requestID, ok := ctx.Value(common.RequestIDContextKey).(string)
	if !ok {
		db.logger.Warn("AddNewRefreshToken: failed to get requestID from context")
		requestID = "unknown"
	}
	db.logger.Debug("AddNewRefreshToken called",
		db.logger.ToString("requestID", requestID),
		db.logger.ToInt("user_id", int(userID)),
	)

	query := `
		INSERT INTO user_session (user_id, session_token, expires_at)
		VALUES ($1, $2, $3)
	`

	_, err := db.conn.Exec(query, userID, refreshToken, expiresAt)
	if err != nil {
		return err
	}

	return nil
}

// Return all refresh tokens for a given user
func (db *DataBase) GetRefreshTokensForUser(ctx context.Context, userID uint) ([]entity.RefreshToken, error) {
	// Log the request ID from context for tracing
	requestID, ok := ctx.Value(common.RequestIDContextKey).(string)
	if !ok {
		db.logger.Warn("GetRefreshTokens: failed to get requestID from context")
		requestID = "unknown"
	}
	db.logger.Debug("GetRefreshTokens called",
		db.logger.ToString("requestID", requestID),
		db.logger.ToInt("user_id", int(userID)),
	)

	query := `
		SELECT user_session_id, session_token, expires_at
		FROM user_session
		WHERE user_id = $1
	`

	rows, err := db.conn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []entity.RefreshToken
	for rows.Next() {
		var token entity.RefreshToken
		if err := rows.Scan(&token.ID, &token.Token, &token.ExpiresAt); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (db *DataBase) RemoveRefreshToken(ctx context.Context, userID uint, refreshToken string) error {
	// Log the request ID from context for tracing
	requestID, ok := ctx.Value(common.RequestIDContextKey).(string)
	if !ok {
		db.logger.Warn("RemoveRefreshToken: failed to get requestID from context")
		requestID = "unknown"
	}
	db.logger.Debug("RemoveRefreshToken called",
		db.logger.ToString("requestID", requestID),
		db.logger.ToInt("user_id", int(userID)),
	)

	query := `
		DELETE FROM user_session
		WHERE user_id = $1 AND session_token = $2
	`

	_, err := db.conn.Exec(query, userID, refreshToken)
	if err != nil {
		return err
	}

	return nil
}

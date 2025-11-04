package postgres

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
)

func (db *DataBase) AddNewRefreshToken(ctx context.Context, userID uint, refreshToken string, expiresAt time.Time) error {
	// Log the request ID from context for tracing
	requestID, ok := ctx.Value(common.RequestIDContextKey).(string)
	if !ok {
		db.logger.Warn("AddNewRefreshToken: failed to get requestID from context")
		requestID = "unknown"
	}
	db.logger.Info("AddNewRefreshToken called",
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

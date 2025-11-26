package postgres

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

func (db *DataBase) AddNewRefreshToken(ctx context.Context, userID uint, refreshToken string, expiresAt time.Time) error {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("AddNewRefreshToken called",
		log.ToInt("user_id", int(userID)),
		log.ToString("refresh_token", refreshToken),
		log.ToString("expires_at", expiresAt.GoString()),
	)

	query := `
		INSERT INTO user_session (user_id, session_token, expires_at)
		VALUES ($1, $2, $3)
	`

	_, err := db.conn.Exec(query, userID, refreshToken, expiresAt)
	if err != nil {
		log.Error("Failed to add new refresh token: " + err.Error())
		return err
	}

	return nil
}

// Return all refresh tokens for a given user
func (db *DataBase) GetRefreshTokensForUser(ctx context.Context, userID uint) ([]entity.RefreshToken, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetRefreshTokensForUser called",
		log.ToInt("user_id", int(userID)),
	)

	query := `
		SELECT user_session_id, session_token, expires_at
		FROM user_session
		WHERE user_id = $1
	`

	rows, err := db.conn.Query(query, userID)
	if err != nil {
		log.Error("Failed to get refresh tokens for user: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var tokens []entity.RefreshToken
	for rows.Next() {
		var token entity.RefreshToken
		if err := rows.Scan(&token.ID, &token.Token, &token.ExpiresAt); err != nil {
			log.Error("Failed to scan refresh token: " + err.Error())
			return nil, err
		}
		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		log.Error("Row iteration error: " + err.Error())
		return nil, err
	}

	return tokens, nil
}

func (db *DataBase) RemoveRefreshToken(ctx context.Context, userID uint, refreshToken string) error {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("RemoveRefreshToken called",
		log.ToInt("user_id", int(userID)),
		log.ToString("refresh_token", refreshToken),
	)

	query := `
		DELETE FROM user_session
		WHERE user_id = $1 AND session_token = $2
	`

	_, err := db.conn.Exec(query, userID, refreshToken)
	if err != nil {
		log.Error("Failed to remove refresh token: " + err.Error())
		return err
	}

	return nil
}

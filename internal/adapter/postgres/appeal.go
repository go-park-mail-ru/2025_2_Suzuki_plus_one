package postgres

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

func (db *DataBase) GetAppealIDsByUserID(ctx context.Context, userID uint) ([]uint, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetAppealIDsByUserID called",
		log.ToInt("user_id", int(userID)),
	)

	var appealIDs []uint
	query := "SELECT user_appeal_id FROM user_appeal WHERE user_id = $1"
	rows, err := db.conn.QueryContext(ctx, query, userID)
	if err != nil {
		log.Error("Failed to query appeal IDs: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var appealID uint
		if err := rows.Scan(&appealID); err != nil {
			log.Error("Failed to scan appeal ID: " + err.Error())
			return nil, err
		}
		appealIDs = append(appealIDs, appealID)
	}
	return appealIDs, nil
}

func (db *DataBase) GetAppealIDsAll(ctx context.Context, tag *string, status *string, limit uint, offset uint) ([]uint, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetAppealIDsAll called")
	if tag == nil {
		log.Warn("tag", log.ToString("tag", "nil"))
	}
	if status == nil {
		log.Warn("status", log.ToString("status", "nil"))
	}

	var appealIDs []uint
	query := `
		SELECT user_appeal_id FROM user_appeal
		WHERE ($1::text IS NULL OR tag = $1::text)
		AND   ($2::text IS NULL OR status = $2::text)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := db.conn.QueryContext(ctx, query, tag, status, limit, offset)
	if err != nil {
		log.Error("Failed to query appeal IDs: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var appealID uint
		if err := rows.Scan(&appealID); err != nil {
			log.Error("Failed to scan appeal ID: " + err.Error())
			return nil, err
		}
		appealIDs = append(appealIDs, appealID)
	}
	return appealIDs, nil
}

func (db *DataBase) GetAppealByID(ctx context.Context, appealID uint) (*entity.Appeal, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetAppealByID called",
		log.ToInt("appeal_id", int(appealID)),
	)

	var appeal entity.Appeal
	query := `
		SELECT user_appeal_id, user_id, tag, status, name, created_at, updated_at
		FROM user_appeal
		WHERE user_appeal_id = $1
	`
	err := db.conn.QueryRowContext(ctx, query, appealID).Scan(
		&appeal.ID,
		&appeal.UserID,
		&appeal.Tag,
		&appeal.Status,
		&appeal.Name,
		&appeal.CreatedAt,
		&appeal.UpdatedAt,
	)
	if err != nil {
		log.Error("Failed to get appeal by ID: " + err.Error())
		return nil, err
	}
	return &appeal, nil
}

func (db *DataBase) GetAppealMessagesIDsByAppealID(ctx context.Context, appealID uint) ([]uint, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetAppealMessagesIDsByAppealID called",
		log.ToInt("appeal_id", int(appealID)),
	)

	var messageIDs []uint
	query := `
		SELECT user_appeal_message_id FROM user_appeal_message
		WHERE user_appeal_id = $1
		ORDER BY created_at DESC
	`
	rows, err := db.conn.QueryContext(ctx, query, appealID)
	if err != nil {
		log.Error("Failed to query appeal message IDs: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var messageID uint
		if err := rows.Scan(&messageID); err != nil {
			log.Error("Failed to scan appeal message ID: " + err.Error())
			return nil, err
		}
		messageIDs = append(messageIDs, messageID)
	}
	return messageIDs, nil
}

func (db *DataBase) GetAppealMessagesByID(ctx context.Context, appealID uint) ([]entity.AppealMessage, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetAppealMessagesByID called",
		log.ToInt("appeal_id", int(appealID)),
	)

	var messages []entity.AppealMessage

	// Retrieve all messages from user
	query := `
		SELECT user_appeal_message_id, is_response, message_content, created_at
		FROM user_appeal_message
		WHERE user_appeal_id = $1
		ORDER BY created_at DESC
	`

	rows, err := db.conn.QueryContext(ctx, query, appealID)
	if err != nil {
		log.Error("Failed to query appeal messages: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var message entity.AppealMessage
		if err := rows.Scan(
			&message.ID,
			&message.IsResponse,
			&message.Message,
			&message.CreatedAt,
		); err != nil {
			log.Error("Failed to scan appeal message: " + err.Error())
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// Create

func (db *DataBase) CreateAppeal(ctx context.Context, userID uint, tag string, name string) (uint, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("CreateAppeal called",
		log.ToInt("user_id", int(userID)),
		log.ToString("tag", tag),
		log.ToString("name", name),
	)

	var appealID uint
	query := `
		INSERT INTO user_appeal (user_id, tag, name, status)
		VALUES ($1, $2, $3, 'open')
		RETURNING user_appeal_id
	`
	err := db.conn.QueryRowContext(ctx, query, userID, tag, name).Scan(&appealID)
	if err != nil {
		log.Error("Failed to create appeal: " + err.Error())
		return 0, err
	}
	return appealID, nil
}

func (db *DataBase) CreateAppealMessage(ctx context.Context, appealID uint, isResponse bool, message string) (uint, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("CreateAppealMessage called",
		log.ToInt("appeal_id", int(appealID)),
		log.ToAny("is_response", isResponse),
	)
	var messageID uint
	query := `
		INSERT INTO user_appeal_message (user_appeal_id, is_response, message_content)
		VALUES ($1, $2, $3)
		RETURNING user_appeal_message_id
	`
	err := db.conn.QueryRowContext(ctx, query, appealID, isResponse, message).Scan(&messageID)
	if err != nil {
		log.Error("Failed to create appeal message: " + err.Error())
		return 0, err
	}
	return messageID, nil
}

// Update

func (db *DataBase) UpdateAppealStatus(ctx context.Context, appealID uint, status string) error {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("UpdateAppealStatus called",
		log.ToInt("appeal_id", int(appealID)),
		log.ToString("status", status),
	)

	query := `
		UPDATE user_appeal
		SET status = $1, updated_at = NOW()
		WHERE user_appeal_id = $2
	`
	_, err := db.conn.ExecContext(ctx, query, status, appealID)
	if err != nil {
		log.Error("Failed to update appeal status: " + err.Error())
		return err
	}
	return nil
}

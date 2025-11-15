package postgres

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

func (db *DataBase) GetAppealIDsByUserID(ctx context.Context, userID uint) ([]uint, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetAppealIDsByUserID called",
		log.ToInt("user_id", int(userID)),
	)

	var appealIDs []uint
	query := "SELECT appeal_id FROM appeal WHERE user_id = $1"
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

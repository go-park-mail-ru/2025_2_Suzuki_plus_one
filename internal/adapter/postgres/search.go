package postgres

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

func (db *DataBase) SearchMedia(ctx context.Context, query string, limit, offset uint) ([]uint, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("SearchMedia called",
		log.ToString("query", query),
	)

	sql_query := `
		SELECT media_id
		FROM media
		WHERE title ILIKE '%' || $1 || '%'
		ORDER BY rating DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.conn.QueryContext(ctx, sql_query, query, limit, offset)
	if err != nil {
		log.Error("SearchMedia query failed", log.ToError(err))
		return nil, err
	}
	defer rows.Close()

	var mediaIDs []uint
	for rows.Next() {
		var mediaID uint
		if err := rows.Scan(&mediaID); err != nil {
			log.Error("SearchMedia row scan failed", log.ToError(err))
			return nil, err
		}
		mediaIDs = append(mediaIDs, mediaID)
	}

	if err := rows.Err(); err != nil {
		log.Error("SearchMedia rows iteration error", log.ToError(err))
		return nil, err
	}

	log.Info("SearchMedia completed successfully", log.ToInt("results_count", len(mediaIDs)))
	return mediaIDs, nil
}


func (db *DataBase) SearchActor(ctx context.Context, query string, limit, offset uint) ([]uint, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("SearchActor called",
		log.ToString("query", query),
	)

	sql_query := `
		SELECT actor_id
		FROM actor
		WHERE name ILIKE '%' || $1 || '%'
		LIMIT $2 OFFSET $3
	`

	rows, err := db.conn.QueryContext(ctx, sql_query, query, limit, offset)
	if err != nil {
		log.Error("SearchActor query failed", log.ToError(err))
		return nil, err
	}
	defer rows.Close()

	var actorIDs []uint
	for rows.Next() {
		var actorID uint
		if err := rows.Scan(&actorID); err != nil {
			log.Error("SearchActor row scan failed", log.ToError(err))
			return nil, err
		}
		actorIDs = append(actorIDs, actorID)
	}

	if err := rows.Err(); err != nil {
		log.Error("SearchActor rows iteration error", log.ToError(err))
		return nil, err
	}

	log.Info("SearchActor completed successfully", log.ToInt("results_count", len(actorIDs)))
	return actorIDs, nil
}
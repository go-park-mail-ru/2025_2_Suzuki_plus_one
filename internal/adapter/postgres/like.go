package postgres

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

func (db *DataBase) GetLike(ctx context.Context, userID uint, mediaID uint) (exists bool, ent entity.Like, err error) {
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetLike called",
		log.ToInt("user_id", int(userID)),
		log.ToInt("media_id", int(mediaID)),
	)

	query := `
	SELECT user_id, media_id, is_dislike
	FROM user_like_media
	WHERE user_id = $1 AND media_id = $2
	`

	var like entity.Like
	err = db.conn.QueryRow(query, userID, mediaID).Scan(
		&like.UserID,
		&like.MediaID,
		&like.IsDislike,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug("No like found",
				log.ToInt("user_id", int(userID)),
				log.ToInt("media_id", int(mediaID)),
			)
			return false, entity.Like{}, nil
		}
		return false, entity.Like{}, err
	}

	log.Debug("Like found",
		log.ToInt("user_id", int(userID)),
		log.ToInt("media_id", int(mediaID)),
	)

	return true, like, nil
}


func (db *DataBase) ToggleLike(ctx context.Context, userID uint, mediaID uint) (isDislike bool, err error) {
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("ToggleLike called",
		log.ToInt("user_id", int(userID)),
		log.ToInt("media_id", int(mediaID)),
	)

	query := `
	INSERT INTO user_like_media (user_id, media_id, is_dislike)
	VALUES ($1, $2, FALSE)
	ON CONFLICT (user_id, media_id)
	DO UPDATE SET is_dislike = NOT user_like_media.is_dislike
	RETURNING is_dislike
	`

	var isDislikeResult bool
	err = db.conn.QueryRow(query, userID, mediaID).Scan(&isDislikeResult)
	if err != nil {
		return false, err
	}

	log.Debug("Like toggled",
		log.ToInt("user_id", int(userID)),
		log.ToInt("media_id", int(mediaID)),
		log.ToAny("is_dislike", isDislikeResult),
	)

	return isDislikeResult, nil
}

func (db *DataBase) DeleteLike(ctx context.Context, userID uint, mediaID uint) error {
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("DeleteLike called",
		log.ToInt("user_id", int(userID)),
		log.ToInt("media_id", int(mediaID)),
	)

	query := `
	DELETE FROM user_like_media
	WHERE user_id = $1 AND media_id = $2
	`

	_, err := db.conn.Exec(query, userID, mediaID)
	if err != nil {
		return err
	}

	log.Debug("Like deleted",
		log.ToInt("user_id", int(userID)),
		log.ToInt("media_id", int(mediaID)),
	)

	return nil
}
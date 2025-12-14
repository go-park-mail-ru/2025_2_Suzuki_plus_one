package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

func (db *DataBase) GetGenreByID(ctx context.Context, genreID uint) (*entity.Genre, error) {
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetGenreByID called",
		log.ToInt("genre_id", int(genreID)),
	)

	var genre entity.Genre
	query := `
	SELECT genre_id, name, description
	FROM genre
	WHERE genre_id = $1
	`
	row := db.conn.QueryRow(query, genreID)
	if err := row.Scan(&genre.ID, &genre.Name, &genre.Description); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("GetGenreByID: genre not found",
				log.ToInt("genre_id", int(genreID)),
			)
			return nil, nil
		}
		log.Error("GetGenreByID: failed to scan genre", log.ToError(err))
		return nil, err
	}
	return &genre, nil
}

func (db *DataBase) GetAllGenreIDs(ctx context.Context) ([]uint, error) {
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetAllGenreIDs called")

	var genreIDs []uint
	query := `
	SELECT genre_id
	FROM genre
	`
	rows, err := db.conn.Query(query)
	if err != nil {
		log.Error("GetAllGenreIDs: failed to execute query", log.ToError(err))
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Error("GetAllGenreIDs: failed to close rows", log.ToError(cerr))
		}
	}()
	for rows.Next() {
		var genreID uint
		if err := rows.Scan(&genreID); err != nil {
			log.Error("GetAllGenreIDs: failed to scan genre ID", log.ToError(err))
			return nil, err
		}
		genreIDs = append(genreIDs, genreID)
	}
	return genreIDs, nil
}

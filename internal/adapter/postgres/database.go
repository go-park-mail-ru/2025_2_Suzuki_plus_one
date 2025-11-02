package postgres

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/jackc/pgx/v5"
)

type DataBase struct {
	logger     logger.Logger
	conn       *pgx.Conn
	connString string
	context    context.Context
}

// Initialize a new Postgres database connection
func NewDataBase(logger logger.Logger, dbUrl string) *DataBase {
	return &DataBase{
		logger:     logger,
		connString: dbUrl,
	}
}

func (db *DataBase) Connect() error {
	db.context = context.Background()
	conn, err := pgx.Connect(db.context, db.connString)
	if err != nil {
		return err
	}
	db.conn = conn
	return nil
}

func (db *DataBase) Close() error {
	return db.conn.Close(db.context)
}

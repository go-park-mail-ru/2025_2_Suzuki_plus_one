package postgres

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DataBase struct {
	logger     logger.Logger
	conn       *sql.DB
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
	conn, err := sql.Open("pgx", db.connString)
	if err != nil {
		return err
	}
	if err := conn.PingContext(db.context); err != nil {
		return err
	}
	db.logger.Info("Database connection established")
	db.conn = conn
	return nil
}

func (db *DataBase) Close() error {
	return db.conn.Close()
}

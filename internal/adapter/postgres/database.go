package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DataBase struct {
	logger      logger.Logger
	conn        *sql.DB
	connString  string
	context     context.Context
	maxOpenCons int
	maxIdleCons int
	maxLifetime time.Duration
}

// Initialize a new Postgres database connection
func NewDataBase(logger logger.Logger, dbUrl string, maxOpenCons, maxIdleCons int, maxLifetime time.Duration) *DataBase {
	return &DataBase{
		logger:      logger,
		connString:  dbUrl,
		maxOpenCons: maxOpenCons,
		maxIdleCons: maxIdleCons,
		maxLifetime: maxLifetime,
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

	// Setup connection pool settings
	db.conn.SetMaxOpenConns(db.maxOpenCons)
	db.conn.SetMaxIdleConns(db.maxIdleCons)
	db.conn.SetConnMaxLifetime(db.maxLifetime)

	// Проверка соединения при старте
	if err := db.conn.PingContext(db.context); err != nil {
		db.logger.Error("Failed to ping database: " + err.Error())
		return err
	}

	return nil
}

func (db *DataBase) Close() error {
	return db.conn.Close()
}

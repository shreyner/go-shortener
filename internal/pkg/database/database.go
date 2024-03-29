// Package database create connection and check connection on error
package database

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

// NewDataBase returning sql.DB witch active connection or error if was error to connection
func NewDataBase(log *zap.Logger, dburi string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dburi)

	if err != nil {
		log.Error("error when connected db", zap.Error(err))
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Error("error when ping to db", zap.Error(err))
		return nil, err
	}

	return db, nil
}

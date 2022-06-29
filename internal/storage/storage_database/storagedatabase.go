package storagedatabase

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/stdlib"
)

type StorageSQL interface {
	PingContext(ctx context.Context) error
}

type storageSQL struct {
	DB *sql.DB
}

func NewStorageSQL(dataBaseDSN string) (
	*storageSQL,
	error,
) {
	db, err := sql.Open("pgx", dataBaseDSN)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &storageSQL{
		DB: db,
	}, nil
}

func (s *storageSQL) PingContext(ctx context.Context) error {
	err := s.DB.PingContext(ctx)

	return err
}

func (s *storageSQL) Close() error {
	err := s.DB.Close()

	return err
}

package storagedatabase

import (
	"context"
	"database/sql"
)

type StorageSQL interface {
	PingContext(ctx context.Context) error
}

type storageSQL struct {
	DB *sql.DB
}

func NewStorageSQL(db *sql.DB) *storageSQL {
	return &storageSQL{
		DB: db,
	}
}

func (s *storageSQL) PingContext(ctx context.Context) error {
	return s.DB.PingContext(ctx)
}

func (s *storageSQL) Close() error {
	err := s.DB.Close()

	return err
}

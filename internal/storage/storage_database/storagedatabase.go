package storagedatabase

import (
	"context"
	"database/sql"
)

// StorageSQL sql storage include db connection
type StorageSQL struct {
	DB *sql.DB
}

// NewStorageSQL create sql store by db connection
func NewStorageSQL(db *sql.DB) *StorageSQL {
	return &StorageSQL{
		DB: db,
	}
}

// PingContext check connection to db
func (s *StorageSQL) PingContext(ctx context.Context) error {
	return s.DB.PingContext(ctx)
}

// Close connection to db
func (s *StorageSQL) Close() error {
	err := s.DB.Close()

	return err
}

package storage

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/shreyner/go-shortener/internal/pkg/database"
	"github.com/shreyner/go-shortener/internal/repositories"
	storagedatabase "github.com/shreyner/go-shortener/internal/storage/storage_database"
	storagefile "github.com/shreyner/go-shortener/internal/storage/storage_file"
	storagememory "github.com/shreyner/go-shortener/internal/storage/storage_memory"
)

const (
	RepositoryTypeFile = iota
	RepositoryTypeDataBase
	RepositoryTypeMemory
)

type Storage struct {
	ShortURL repositories.ShortURLRepository

	ping  func(context.Context) error
	close func() error
}

func NewStorage(log *zap.Logger, fileStoragePath string, dataBaseDSN string) (*Storage, error) {
	var repositoryType int

	if fileStoragePath != "" {
		repositoryType = RepositoryTypeFile
	} else if dataBaseDSN != "" {
		repositoryType = RepositoryTypeDataBase
	} else {
		repositoryType = RepositoryTypeMemory
	}

	if repositoryType == RepositoryTypeFile {
		log.Info("Init file storage")
		shorterFileRepository, err := storagefile.NewShortURLStore(fileStoragePath)

		if err != nil {
			return nil, fmt.Errorf("storage error when initialize file: %w", err)
		}

		return &Storage{
			ShortURL: shorterFileRepository,

			ping:  func(_ context.Context) error { return nil },
			close: shorterFileRepository.Close,
		}, nil
	}

	if repositoryType == RepositoryTypeDataBase {
		log.Info("Init database storage")
		log.Info("Connect to database ...")
		db, err := database.NewDataBase(log, dataBaseDSN)

		if err != nil {
			return nil, fmt.Errorf("storage error when initialize connection to db: %w", err)
		}

		storeDB := storagedatabase.NewStorageSQL(db)

		log.Info("Success connected database")

		log.Info("Check and create database...")
		if err := storeDB.CheckAndCreateSchema(); err != nil {
			return nil, fmt.Errorf("storage error when create schema in db: %w", err)
		}
		log.Info("Finish check or created...")

		shortURLStorage, err := storagedatabase.NewShortURLStore(log, db)
		if err != nil {
			return nil, fmt.Errorf("storage error when initilizing shortURLStorage: %w", err)
		}

		return &Storage{
			ShortURL: shortURLStorage,

			ping: storeDB.PingContext,
			close: func() error {
				log.Info("Close database connection")
				if err := storeDB.Close(); err != nil {
					log.Error("error to close connection db", zap.Error(err))

					return err
				}

				return nil
			},
		}, nil
	}

	log.Info("Init memory storage")
	return &Storage{
		ShortURL: storagememory.NewShortURLStore(),

		ping:  func(_ context.Context) error { return nil },
		close: func() error { return nil },
	}, nil

}

func (s *Storage) PingContext(ctx context.Context) error {
	return s.ping(ctx)
}

func (s *Storage) Close() error {
	return s.close()
}

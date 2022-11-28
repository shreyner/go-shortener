package storagefile

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"go.uber.org/zap"

	"github.com/shreyner/go-shortener/internal/core"
	"github.com/shreyner/go-shortener/internal/repositories"
)

var (
	_ repositories.ShortURLRepository = (*shortURLRepository)(nil)
)

type shortURLRepository struct {
	decoder    *json.Decoder
	encoder    *json.Encoder
	file       *os.File
	mutex      *sync.RWMutex
	log        *zap.Logger
	pathToFile string
}

// NewShortURLStore create file store
func NewShortURLStore(log *zap.Logger, fileStoragePath string) (*shortURLRepository, error) {
	file, err := os.OpenFile(fileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		return nil, err
	}

	return &shortURLRepository{
		log:        log,
		pathToFile: fileStoragePath,
		mutex:      &sync.RWMutex{},
		file:       file,
		decoder:    json.NewDecoder(file),
		encoder:    json.NewEncoder(file),
	}, nil
}

// Add Добавить короткую ссылку в store
func (s *shortURLRepository) Add(_ context.Context, shortURL *core.ShortURL) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.encoder.Encode(&shortURL)
}

// GetByID Получить короткую ссылку по идентификатору
func (s *shortURLRepository) GetByID(ctx context.Context, id string) (*core.ShortURL, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	file, err := os.OpenFile(s.pathToFile, os.O_RDONLY|os.O_CREATE, 0644)

	if err != nil {
		s.log.Error("error open file for read", zap.Error(err))
		return nil, false
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for decoder.More() {
		var shortURL core.ShortURL
		err := decoder.Decode(&shortURL)

		if err != nil {
			s.log.Error("error read shorted json", zap.Error(err))
			return nil, false
		}

		if shortURL.ID == id {
			return &shortURL, true
		}
	}

	return nil, false
}

// AllByUserID получить все ссылки по идентификатору пользователя
func (s *shortURLRepository) AllByUserID(_ context.Context, id string) ([]*core.ShortURL, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	file, err := os.OpenFile(s.pathToFile, os.O_RDONLY|os.O_CREATE, 0644)

	if err != nil {
		s.log.Error("error open file for read", zap.Error(err))
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var result []*core.ShortURL

	for decoder.More() {
		var shortURL core.ShortURL
		err := decoder.Decode(&shortURL)

		if err != nil {
			s.log.Error("error read shorted json", zap.Error(err))
			return nil, err
		}

		if shortURL.UserID != "" && shortURL.UserID == id {
			result = append(result, &shortURL)
		}
	}

	return result, nil
}

// CreateBatch Добавление ссылок пачкой
func (s *shortURLRepository) CreateBatch(_ context.Context, shortURLs *[]*core.ShortURL) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.encoder.Encode(shortURLs)
}

// Close Метод для корректного закрытия store
func (s *shortURLRepository) Close() error {
	return s.file.Close()
}

// DeleteURLsUserByIds Удаление пачкой коротких ссылок от имени пользователя
func (s *shortURLRepository) DeleteURLsUserByIds(_ context.Context, userID string, ids []string) error {
	return nil
}

// GetStats return stats
func (s *shortURLRepository) GetStats(_ context.Context) (*core.ShortStats, error) {
	return &core.ShortStats{}, nil
}

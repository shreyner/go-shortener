package storagefile

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/shreyner/go-shortener/internal/core"
	"github.com/shreyner/go-shortener/internal/repositories"
)

var (
	_ repositories.ShortURLRepository = (*shortURLRepository)(nil)
)

type shortURLRepository struct {
	pathToFile string
	file       *os.File
	decoder    *json.Decoder
	encoder    *json.Encoder

	mutex *sync.RWMutex
}

func NewShortURLStore(fileStoragePath string) (*shortURLRepository, error) {
	file, err := os.OpenFile(fileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		return nil, err
	}

	return &shortURLRepository{
		pathToFile: fileStoragePath,
		mutex:      &sync.RWMutex{},
		file:       file,
		decoder:    json.NewDecoder(file),
		encoder:    json.NewEncoder(file),
	}, nil
}

func (s *shortURLRepository) Add(shortURL *core.ShortURL) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.encoder.Encode(&shortURL)
}

func (s *shortURLRepository) GetByID(id string) (*core.ShortURL, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	file, err := os.OpenFile(s.pathToFile, os.O_RDONLY|os.O_CREATE, 0644)

	if err != nil {
		log.Printf("Error open file for read: %s", err)
		return nil, false
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for decoder.More() {
		var shortURL core.ShortURL
		err := decoder.Decode(&shortURL)

		if err != nil {
			log.Printf("Error read shorted json:%s\n", err)
			return nil, false
		}

		if shortURL.ID == id {
			return &shortURL, true
		}
	}

	return nil, false
}

func (s *shortURLRepository) AllByUserID(id string) ([]*core.ShortURL, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	file, err := os.OpenFile(s.pathToFile, os.O_RDONLY|os.O_CREATE, 0644)

	if err != nil {
		log.Printf("Error open file for read: %s", err)
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var result []*core.ShortURL

	for decoder.More() {
		var shortURL core.ShortURL
		err := decoder.Decode(&shortURL)

		if err != nil {
			log.Printf("Error read shorted json:%s\n", err)
			return nil, err
		}

		if shortURL.UserID != "" && shortURL.UserID == id {
			result = append(result, &shortURL)
		}
	}

	return result, nil
}

func (s *shortURLRepository) CreateBatchWithContext(_ context.Context, shortURLs *[]*core.ShortURL) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.encoder.Encode(shortURLs)
}

func (s *shortURLRepository) Close() error {
	return s.file.Close()
}

func (s *shortURLRepository) DeleteURLsUserByIds(userID string, ids []string) error {
	return nil
}

// Package storagememory хранилище в памяти
package storagememory

import (
	"context"
	"sync"

	"github.com/shreyner/go-shortener/internal/core"
	"github.com/shreyner/go-shortener/internal/repositories"
)

var (
	_ repositories.ShortURLRepository = (*shortURLRepository)(nil)
)

type shortURLRepository struct {
	store map[string]*core.ShortURL
	mutex *sync.RWMutex
}

func NewShortURLStore() *shortURLRepository {
	return &shortURLRepository{
		store: map[string]*core.ShortURL{},
		mutex: &sync.RWMutex{},
	}
}

// Add Добавить короткую ссылку в store
func (s *shortURLRepository) Add(_ context.Context, shortURL *core.ShortURL) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.store[shortURL.ID] = shortURL

	return nil
}

// GetByID Получить короткую ссылку по идентификатору
func (s *shortURLRepository) GetByID(_ context.Context, id string) (*core.ShortURL, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	shortURL, ok := s.store[id]

	return shortURL, ok
}

// AllByUserID получить все ссылки по идентификатору пользователя
func (s *shortURLRepository) AllByUserID(_ context.Context, id string) ([]*core.ShortURL, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var result []*core.ShortURL

	for _, shortURL := range s.store {
		if shortURL.UserID != "" && shortURL.UserID == id {
			result = append(result, shortURL)
		}
	}

	return result, nil
}

// CreateBatch Добавление ссылок пачкой
func (s *shortURLRepository) CreateBatch(_ context.Context, shortURLs *[]*core.ShortURL) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, v := range *shortURLs {
		s.store[v.ID] = v
	}

	return nil
}

// DeleteURLsUserByIds Удаление пачкой коротких ссылок от имени пользователя
func (s *shortURLRepository) DeleteURLsUserByIds(_ context.Context, userID string, ids []string) error {
	return nil
}

package storagememory

import (
	"sync"

	"github.com/shreyner/go-shortener/internal/core"
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

func (s *shortURLRepository) Add(shortURL *core.ShortURL) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.store[shortURL.ID] = shortURL

	return nil
}

func (s *shortURLRepository) GetByID(id string) (*core.ShortURL, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	shortURL, ok := s.store[id]

	return shortURL, ok
}

func (s *shortURLRepository) AllByUserID(id string) ([]*core.ShortURL, error) {
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

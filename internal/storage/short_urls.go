package storage

import (
	sync "sync"

	"github.com/shreyner/go-shortener/internal/core"
)

type ShortURLRepository struct {
	store map[string]core.ShortURL
	mutex *sync.RWMutex
}

func NewShortURLStore() *ShortURLRepository {
	return &ShortURLRepository{
		store: map[string]core.ShortURL{},
		mutex: &sync.RWMutex{},
	}
}

func (s *ShortURLRepository) Add(shortURL core.ShortURL) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.store[shortURL.ID] = shortURL

	return nil
}

func (s *ShortURLRepository) GetByID(id string) (core.ShortURL, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	shortURL, ok := s.store[id]

	return shortURL, ok
}

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

func (s *shortURLRepository) CreateBatchWithContext(_ context.Context, shortURLs *[]*core.ShortURL) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, v := range *shortURLs {
		s.store[v.ID] = v
	}

	return nil
}

func (s *shortURLRepository) DeleteURLsUserByIds(userID string, ids []string) error {
	return nil
}

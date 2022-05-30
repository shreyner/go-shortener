package storage

import "github.com/shreyner/go-shortener/internal/core"

type ShortURLRepository struct {
	store map[string]core.ShortURL
}

func NewShortURLStore() *ShortURLRepository {
	return &ShortURLRepository{
		store: map[string]core.ShortURL{},
	}
}

func (s *ShortURLRepository) Add(shortURL core.ShortURL) error {
	s.store[shortURL.ID] = shortURL

	return nil
}

func (s *ShortURLRepository) GetByID(id string) (core.ShortURL, bool) {
	shortURL, ok := s.store[id]
	return shortURL, ok
}

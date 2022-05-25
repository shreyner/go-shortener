package storage

import "github.com/shreyner/go-shortener/internal/core"

type ShortUrlRepository struct {
	store map[string]core.ShortUrl
}

func NewShortUrlStore() *ShortUrlRepository {
	return &ShortUrlRepository{
		store: map[string]core.ShortUrl{},
	}
}

func (s *ShortUrlRepository) Add(shortedUrl core.ShortUrl) error {
	s.store[shortedUrl.Id] = shortedUrl

	return nil
}

func (s *ShortUrlRepository) GetById(id string) (core.ShortUrl, bool) {
	shortUrl, ok := s.store[id]
	return shortUrl, ok
}

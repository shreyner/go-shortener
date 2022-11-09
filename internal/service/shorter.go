package service

import (
	"context"

	"github.com/shreyner/go-shortener/internal/core"
	rand "github.com/shreyner/go-shortener/internal/pkg/random"
	"github.com/shreyner/go-shortener/internal/repositories"
)

var (
	lengthShortID = 10
)

// Shorter service include business logic for work with short URLs
type Shorter struct {
	shorterRepository repositories.ShortURLRepository
}

// NewShorter create service
func NewShorter(shorterRepository repositories.ShortURLRepository) *Shorter {
	return &Shorter{
		shorterRepository: shorterRepository,
	}
}

// Create new short url by user
func (s *Shorter) Create(ctx context.Context, userID, url string) (*core.ShortURL, error) {
	id := generateURLID()
	shortURL := &core.ShortURL{ID: id, URL: url, UserID: userID}

	err := s.shorterRepository.Add(ctx, shortURL)

	if err != nil {
		return nil, err
	}

	return shortURL, nil
}

// CreateBatch more URLs by user
func (s *Shorter) CreateBatch(ctx context.Context, shortURLs *[]*core.ShortURL) error {
	for _, v := range *shortURLs {
		v.ID = generateURLID()
	}

	if err := s.shorterRepository.CreateBatch(ctx, shortURLs); err != nil {
		return err
	}

	return nil
}

// GetByID find by short URL and return original url or error with not found
func (s *Shorter) GetByID(ctx context.Context, id string) (*core.ShortURL, bool) {
	return s.shorterRepository.GetByID(ctx, id)
}

// AllByUser return was created user
func (s *Shorter) AllByUser(ctx context.Context, id string) ([]*core.ShortURL, error) {
	return s.shorterRepository.AllByUserID(ctx, id)
}

func generateURLID() string {
	return rand.RandSeq(lengthShortID)
}

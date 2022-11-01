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

type Shorter struct {
	shorterRepository repositories.ShortURLRepository
}

func NewShorter(shorterRepository repositories.ShortURLRepository) *Shorter {
	return &Shorter{
		shorterRepository: shorterRepository,
	}
}

func (s *Shorter) Create(ctx context.Context, userID, url string) (*core.ShortURL, error) {
	id := generateURLID()
	shortURL := &core.ShortURL{ID: id, URL: url, UserID: userID}

	err := s.shorterRepository.Add(ctx, shortURL)

	if err != nil {
		return nil, err
	}

	return shortURL, nil
}

func (s *Shorter) CreateBatch(ctx context.Context, shortURLs *[]*core.ShortURL) error {
	for _, v := range *shortURLs {
		v.ID = generateURLID()
	}

	if err := s.shorterRepository.CreateBatch(ctx, shortURLs); err != nil {
		return err
	}

	return nil
}

func (s *Shorter) GetByID(ctx context.Context, id string) (*core.ShortURL, bool) {
	return s.shorterRepository.GetByID(ctx, id)
}

func (s *Shorter) AllByUser(ctx context.Context, id string) ([]*core.ShortURL, error) {
	return s.shorterRepository.AllByUserID(ctx, id)
}

func generateURLID() string {
	return rand.RandSeq(lengthShortID)
}

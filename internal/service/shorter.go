package service

import (
	"context"
	"github.com/shreyner/go-shortener/internal/core"
	rand "github.com/shreyner/go-shortener/internal/pkg/random"
	"github.com/shreyner/go-shortener/internal/repositories"
)

var (
	lengthShortID = 4
)

type Shorter struct {
	shorterRepository repositories.ShortURLRepository
}

func NewShorter(shorterRepository repositories.ShortURLRepository) *Shorter {
	return &Shorter{shorterRepository: shorterRepository}
}

func (s *Shorter) Create(userID, url string) (*core.ShortURL, error) {
	id := generateURLID()
	shortURL := core.ShortURL{ID: id, URL: url, UserID: userID}

	err := s.shorterRepository.Add(&shortURL)

	if err != nil {
		return nil, err
	}

	return &shortURL, nil
}

func (s *Shorter) CreateBatchWithContext(ctx context.Context, shortURLs *[]*core.ShortURL) error {
	for _, v := range *shortURLs {
		v.ID = generateURLID()
	}

	if err := s.shorterRepository.CreateBatchWithContext(ctx, shortURLs); err != nil {
		return err
	}

	return nil
}

func (s *Shorter) GetByID(id string) (*core.ShortURL, bool) {
	return s.shorterRepository.GetByID(id)
}

func (s *Shorter) AllByUser(id string) ([]*core.ShortURL, error) {
	return s.shorterRepository.AllByUserID(id)
}

func generateURLID() string {
	return rand.RandSeq(lengthShortID)
}

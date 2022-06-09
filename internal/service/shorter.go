package service

import (
	"math/rand"

	"github.com/shreyner/go-shortener/internal/core"
)

type ShortURLRepository interface {
	Add(shortedURL core.ShortURL) error
	GetByID(id string) (core.ShortURL, bool)
}

type Shorter struct {
	shorterRepository ShortURLRepository
}

func NewShorter(shorterRepository ShortURLRepository) *Shorter {
	return &Shorter{shorterRepository: shorterRepository}
}

func (s *Shorter) Create(url string) (core.ShortURL, error) {
	id := generateURLID()
	shortURL := core.ShortURL{ID: id, URL: url}

	err := s.shorterRepository.Add(shortURL)

	if err != nil {
		return core.ShortURL{}, err
	}

	return shortURL, nil
}

func (s *Shorter) GetByID(id string) (core.ShortURL, bool) {
	return s.shorterRepository.GetByID(id)
}

func generateURLID() string {
	return randSeq(4)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

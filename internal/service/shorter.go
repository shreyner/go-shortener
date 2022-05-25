package service

import (
	"github.com/shreyner/go-shortener/internal/core"
	"math/rand"
)

type ShortUrlRepository interface {
	Add(shortedUrl core.ShortUrl) error
	GetById(id string) (core.ShortUrl, bool)
}

type Shorter struct {
	shorterRepository ShortUrlRepository
}

func NewShorter(shorterRepository ShortUrlRepository) *Shorter {
	return &Shorter{shorterRepository: shorterRepository}
}

func (s *Shorter) Create(url string) (core.ShortUrl, error) {
	id := generateUrlId()
	shortUrl := core.ShortUrl{Id: id, Url: url}

	err := s.shorterRepository.Add(shortUrl)

	if err != nil {
		return core.ShortUrl{}, err
	}

	return shortUrl, nil
}

func (s *Shorter) GetById(id string) (core.ShortUrl, bool) {
	return s.shorterRepository.GetById(id)
}

func generateUrlId() string {
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

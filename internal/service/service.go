package service

import "github.com/shreyner/go-shortener/internal/repositories"

type Services struct {
	ShorterService *Shorter
}

func NewService(shorterRepository repositories.ShortURLRepository) *Services {
	return &Services{
		ShorterService: NewShorter(shorterRepository),
	}
}

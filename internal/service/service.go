package service

import "github.com/shreyner/go-shortener/internal/repositories"

// Services include all services
type Services struct {
	ShorterService *Shorter
}

// NewService return one struct with all services
func NewService(shorterRepository repositories.ShortURLRepository) *Services {
	return &Services{
		ShorterService: NewShorter(shorterRepository),
	}
}

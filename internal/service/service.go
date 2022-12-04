package service

import (
	"go.uber.org/zap"

	"github.com/shreyner/go-shortener/internal/repositories"
)

// Services include all services
type Services struct {
	ShorterService *Shorter
	AuthService    *AuthService
}

// NewService return one struct with all services
func NewService(
	log *zap.Logger,
	shorterRepository repositories.ShortURLRepository,
	signKey []byte,
) (*Services, error) {
	authService, err := NewAuthService(log, signKey)

	if err != nil {
		return nil, err
	}

	services := Services{
		ShorterService: NewShorter(shorterRepository),
		AuthService:    authService,
	}

	return &services, nil
}

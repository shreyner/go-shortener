package repositories

import "github.com/shreyner/go-shortener/internal/core"

type ShortURLRepository interface {
	Add(shortedURL *core.ShortURL) error
	GetByID(id string) (*core.ShortURL, bool)
}

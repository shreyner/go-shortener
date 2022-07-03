package repositories

import (
	"context"

	"github.com/shreyner/go-shortener/internal/core"
)

type ShortURLRepository interface {
	Add(shortedURL *core.ShortURL) error
	GetByID(id string) (*core.ShortURL, bool)
	AllByUserID(id string) ([]*core.ShortURL, error)
	CreateBatchWithContext(ctx context.Context, shortURLs *[]*core.ShortURL) error
}

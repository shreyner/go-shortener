package repositories

import (
	"context"

	"github.com/shreyner/go-shortener/internal/core"
)

// ShortURLRepository base contract for all repositories
type ShortURLRepository interface {
	Add(ctx context.Context, shortedURL *core.ShortURL) error
	GetByID(ctx context.Context, id string) (*core.ShortURL, bool)
	AllByUserID(ctx context.Context, id string) ([]*core.ShortURL, error)
	CreateBatch(ctx context.Context, shortURLs *[]*core.ShortURL) error
	DeleteURLsUserByIds(ctx context.Context, userID string, ids []string) error
	GetStats(ctx context.Context) (*core.ShortStats, error)
}

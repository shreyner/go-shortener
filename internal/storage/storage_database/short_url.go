package storagedatabase

import (
	"context"
	"database/sql"
	"github.com/shreyner/go-shortener/internal/core"
	"time"
)

type shortURLRepository struct {
	db *sql.DB
}

func NewShortURLStore(db *sql.DB) (*shortURLRepository, error) {
	return &shortURLRepository{
		db: db,
	}, nil
}

func (s *shortURLRepository) Add(shortURL *core.ShortURL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		`insert into short_url (id, url, user_id) values ($1, $2, $3);`,
		shortURL.ID,
		shortURL.URL,
		shortURL.UserID,
	)

	return err
}

func (s *shortURLRepository) GetByID(id string) (*core.ShortURL, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var shortURL core.ShortURL

	row := s.db.QueryRowContext(
		ctx,
		`select id, url, user_id from short_url where id = $1`,
		id,
	)

	if row.Err() != nil {
		return nil, false
	}

	if err := row.Scan(&shortURL.ID, &shortURL.URL, &shortURL.UserID); err != nil {
		return nil, false
	}

	return &shortURL, true
}

func (s *shortURLRepository) AllByUserID(id string) ([]*core.ShortURL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(
		ctx,
		`select id, url, user_id from short_url where user_id = $1`,
		id,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var shortURLs []*core.ShortURL

	for rows.Next() {
		shortURL := core.ShortURL{}

		if err := rows.Scan(&shortURL.ID, &shortURL.URL, &shortURL.UserID); err != nil {
			return nil, err
		}

		shortURLs = append(shortURLs, &shortURL)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return shortURLs, nil
}

package storagedatabase

import (
	"context"
	"database/sql"
	"github.com/shreyner/go-shortener/internal/core"
	"time"
)

type shortURLRepository struct {
	db         *sql.DB
	insertStmt *sql.Stmt
}

func NewShortURLStore(db *sql.DB) (*shortURLRepository, error) {
	insertStmt, err := db.Prepare("insert into short_url (id, url, user_id, correlation_id) values ($1, $2, $3, $4);")

	if err != err {
		return nil, err
	}

	return &shortURLRepository{
		db:         db,
		insertStmt: insertStmt,
	}, nil
}

func (s *shortURLRepository) Add(shortURL *core.ShortURL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := s.db.QueryRowContext(
		ctx,
		`insert into short_url (id, url, user_id) values ($1, $2, $3) on conflict (url) do update set url=excluded.url returning id;`,
		shortURL.ID,
		shortURL.URL,
		shortURL.UserID,
	)

	if result.Err() != nil {
		return result.Err()
	}

	var resultID string
	if err := result.Scan(&resultID); err != nil {
		return err
	}

	if resultID != shortURL.ID {
		return NewShortURLCreateConflictError(resultID)
	}

	return nil
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

func (s *shortURLRepository) CreateBatchWithContext(ctx context.Context, shortURLs *[]*core.ShortURL) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}
	defer tx.Rollback()

	txStmt := tx.StmtContext(ctx, s.insertStmt)

	defer txStmt.Close()

	for _, v := range *shortURLs {
		if _, err := txStmt.ExecContext(ctx, v.ID, v.URL, v.UserID, v.CorrelationID); err != nil {
			return err
		}
	}

	return tx.Commit()
}
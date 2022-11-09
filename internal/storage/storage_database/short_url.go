package storagedatabase

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	"github.com/shreyner/go-shortener/internal/core"
	"github.com/shreyner/go-shortener/internal/repositories"
	storeerrors "github.com/shreyner/go-shortener/internal/storage/store_errors"
)

var (
	_ repositories.ShortURLRepository = (*shortURLRepository)(nil)
)

type shortURLRepository struct {
	log        *zap.Logger
	db         *sql.DB
	insertStmt *sql.Stmt
}

// NewShortURLStore create sql store
func NewShortURLStore(log *zap.Logger, db *sql.DB) (*shortURLRepository, error) {
	insertStmt, err := db.Prepare("insert into short_url (id, url, user_id, correlation_id) values ($1, $2, $3, $4);")

	if err != nil {
		return nil, err
	}

	return &shortURLRepository{
		log:        log,
		db:         db,
		insertStmt: insertStmt,
	}, nil
}

// Add Добавить короткую ссылку в store
func (s *shortURLRepository) Add(ctx context.Context, shortURL *core.ShortURL) error {
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
		return storeerrors.NewShortURLCreateConflictError(resultID)
	}

	return nil
}

// GetByID Получить короткую ссылку по идентификатору
func (s *shortURLRepository) GetByID(ctx context.Context, id string) (*core.ShortURL, bool) {
	var shortURL core.ShortURL

	row := s.db.QueryRowContext(
		ctx,
		`select id, url, user_id, deleted from short_url where id = $1`,
		id,
	)

	if row.Err() != nil {
		return nil, false
	}

	if err := row.Scan(&shortURL.ID, &shortURL.URL, &shortURL.UserID, &shortURL.IsDeleted); err != nil {
		return nil, false
	}

	return &shortURL, true
}

// AllByUserID получить все ссылки по идентификатору пользователя
func (s *shortURLRepository) AllByUserID(ctx context.Context, id string) ([]*core.ShortURL, error) {
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

// CreateBatch Добавление ссылок пачкой
func (s *shortURLRepository) CreateBatch(ctx context.Context, shortURLs *[]*core.ShortURL) error {
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			s.log.Error("shortURL Rollback error", zap.Error(err))
		}
	}(tx)

	txStmt := tx.StmtContext(ctx, s.insertStmt)

	defer txStmt.Close()

	for _, v := range *shortURLs {
		if _, err := txStmt.ExecContext(ctx, v.ID, v.URL, v.UserID, v.CorrelationID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// DeleteURLsUserByIds Удаление пачкой коротких ссылок от имени пользователя
func (s *shortURLRepository) DeleteURLsUserByIds(ctx context.Context, userID string, ids []string) error {
	s.log.Info("Was deleted", zap.String("userID", userID), zap.Strings("ids", ids))

	_, err := s.db.ExecContext(
		ctx,
		`update short_url set deleted = true where user_id = $1 and id = any ($2);`,
		userID,
		ids,
	)

	return err
}

package storeerrors

import (
	"fmt"
)

// ShortURLCreateConflictError ошибка конфликта создания короткой ссылки. Содержит идентификатор ранее созданной ссылки
type ShortURLCreateConflictError struct {
	OriginID string
}

func (s *ShortURLCreateConflictError) Error() string {
	return fmt.Sprintf("this url id: %s was created", s.OriginID)
}

func (s *ShortURLCreateConflictError) Unwrap() error {
	return nil
}

func NewShortURLCreateConflictError(originID string) error {
	return &ShortURLCreateConflictError{
		OriginID: originID,
	}
}

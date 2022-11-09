package storeerrors

import (
	"fmt"
)

// ShortURLCreateConflictError conflict created shorter and has ID was create shorter
type ShortURLCreateConflictError struct {
	OriginID string
}

// Error return string about error
func (s *ShortURLCreateConflictError) Error() string {
	return fmt.Sprintf("this url id: %s was created", s.OriginID)
}

// Unwrap for unwrapped error. Not implemented
func (s *ShortURLCreateConflictError) Unwrap() error {
	return nil
}

// NewShortURLCreateConflictError constructor error
func NewShortURLCreateConflictError(originID string) error {
	return &ShortURLCreateConflictError{
		OriginID: originID,
	}
}

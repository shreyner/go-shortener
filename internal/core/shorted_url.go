package core

import "database/sql"

// ShortURL models for short urls
type ShortURL struct {
	ID            string         `json:"id"`
	URL           string         `json:"url"`
	CorrelationID string         `json:"correlation_id,omitempty"`
	UserID        sql.NullString `json:"userId,omitempty"`
	IsDeleted     bool           `json:"isDeleted"`
}

// ShortStats models with stats service
type ShortStats struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}

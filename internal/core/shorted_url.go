package core

// ShortURL models for short urls
type ShortURL struct {
	ID            string `json:"id"`
	URL           string `json:"url"`
	UserID        string `json:"userId,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
	IsDeleted     bool   `json:"isDeleted"`
}

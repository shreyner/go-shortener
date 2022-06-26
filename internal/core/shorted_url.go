package core

type ShortURL struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	UserID string `json:"userId,omitempty"`
}

package storage

type Storage struct {
	ShortUrlRepository *ShortURLRepository
}

func NewStorage() *Storage {
	return &Storage{
		ShortUrlRepository: NewShortURLStore(),
	}
}

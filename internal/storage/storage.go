package storage

type Storage struct {
	ShortUrlRepository *ShortUrlRepository
}

func NewStorage() *Storage {
	return &Storage{
		ShortUrlRepository: NewShortUrlStore(),
	}
}

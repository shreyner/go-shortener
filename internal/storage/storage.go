package storage

type Storage struct {
	ShortURLRepository *ShortURLRepository
}

func NewStorage() *Storage {
	return &Storage{
		ShortURLRepository: NewShortURLStore(),
	}
}

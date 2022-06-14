package storageMemory

type MemoryStorage struct {
	ShortURLRepository *shortURLRepository
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		ShortURLRepository: NewShortURLStore(),
	}
}

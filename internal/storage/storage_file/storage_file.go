package storageFile

type FileStorage struct {
	ShortURLRepository *shortURLRepository
}

func NewFileStorage(fileStoragePath string) (*FileStorage, error) {
	shortURLRepository, err := NewShortURLStore(fileStoragePath)

	if err != nil {
		return nil, err
	}

	return &FileStorage{
		ShortURLRepository: shortURLRepository,
	}, nil
}

func (f *FileStorage) Close() error {
	return f.ShortURLRepository.Close()
}

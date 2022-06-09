package service

type Services struct {
	ShorterService *Shorter
}

func NewService(shorterRepository ShortURLRepository) *Services {
	return &Services{
		ShorterService: NewShorter(shorterRepository),
	}
}

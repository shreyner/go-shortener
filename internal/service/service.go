package service

type Services struct {
	ShorterService *Shorter
}

func NewService(shorterRepository ShortUrlRepository) *Services {
	return &Services{
		ShorterService: NewShorter(shorterRepository),
	}
}

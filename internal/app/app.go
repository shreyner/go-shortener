package app

import (
	"log"

	"github.com/caarlos0/env/v6"

	"github.com/shreyner/go-shortener/internal/handlers"
	"github.com/shreyner/go-shortener/internal/server"
	"github.com/shreyner/go-shortener/internal/service"
	"github.com/shreyner/go-shortener/internal/storage/storage_file"
	"github.com/shreyner/go-shortener/internal/storage/storage_memory"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func NewApp() {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	// TODO: Need refactoring
	if cfg.FileStoragePath != "" {
		storage, err := storagefile.NewFileStorage(cfg.FileStoragePath)

		if err != nil {
			log.Fatal(err)
			return
		}

		defer storage.Close()

		services := service.NewService(storage.ShortURLRepository)

		router := handlers.NewRouter(cfg.BaseURL, services.ShorterService)
		serv := server.NewServer(cfg.ServerAddress, router)

		serv.Start()

	} else {
		storage := storagememory.NewMemoryStorage()

		services := service.NewService(storage.ShortURLRepository)

		router := handlers.NewRouter(cfg.BaseURL, services.ShorterService)
		serv := server.NewServer(cfg.ServerAddress, router)

		serv.Start()
	}
}

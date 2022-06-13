package app

import (
	"log"

	"github.com/caarlos0/env/v6"

	"github.com/shreyner/go-shortener/internal/handlers"
	"github.com/shreyner/go-shortener/internal/server"
	"github.com/shreyner/go-shortener/internal/service"
	"github.com/shreyner/go-shortener/internal/storage"
)

var listerServerAddress = ":8080"

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func NewApp() {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	storages := storage.NewStorage()
	services := service.NewService(storages.ShortURLRepository)

	router := handlers.NewRouter(cfg.BaseURL, services.ShorterService)
	serv := server.NewServer(cfg.ServerAddress, router)

	serv.Start()
}

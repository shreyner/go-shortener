package app

import (
	"github.com/shreyner/go-shortener/internal/handlers"
	"github.com/shreyner/go-shortener/internal/server"
	"github.com/shreyner/go-shortener/internal/service"
	"github.com/shreyner/go-shortener/internal/storage"
)

var listerServerAddress = ":8080"

func NewApp() {
	storages := storage.NewStorage()
	services := service.NewService(storages.ShortUrlRepository)

	router := handlers.NewRouter(services.ShorterService)
	serv := server.NewServer(listerServerAddress, router)

	serv.Start()
}

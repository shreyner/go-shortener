package app

import (
	"context"
	"github.com/shreyner/go-shortener/internal/pkg/fans"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/shreyner/go-shortener/internal/handlers"
	"github.com/shreyner/go-shortener/internal/server"
	"github.com/shreyner/go-shortener/internal/service"
	"github.com/shreyner/go-shortener/internal/storage"
)

func NewApp(
	log *zap.Logger,
	serverAddress string,
	baseURL string,
	fileStoragePath string,
	dataBaseDSN string,
) {
	log.Info("Start app...")

	log.Info("Create storage...")
	store, err := storage.NewStorage(log, fileStoragePath, dataBaseDSN)

	if err != nil {
		log.Error("", zap.Error(err))
		os.Exit(1)
		return
	}

	defer store.Close()

	services := service.NewService(store.ShortURL)

	fansShortService := fans.NewFansShortService(store.ShortURL, 4)

	r := handlers.NewRouter(log, baseURL, services.ShorterService, store.ShortURL, fansShortService, store)
	serv := server.NewServer(log, serverAddress, r)

	serv.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case x := <-interrupt:
		log.Info("Received a signal.", zap.String("signal", x.String()))
	case err := <-serv.Notify():
		log.Error("Received an error from the start rest api server", zap.Error(err))
	}

	log.Info("Stopping server...")

	if err := serv.Stop(context.Background()); err != nil {
		log.Error("Got an error while stopping th rest api server", zap.Error(err))
	}

	fansShortService.Close()

	log.Info("The app is calling the last defers and will be stopped.")
}

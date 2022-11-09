package main

import (
	"flag"
	logStd "log"

	"go.uber.org/zap"

	"github.com/shreyner/go-shortener/internal/app"
	"github.com/shreyner/go-shortener/internal/config"
	"github.com/shreyner/go-shortener/internal/pkg/logger"
)

func main() {
	var cfg config.Config

	if err := cfg.Parse(); err != nil {
		logStd.Fatal("error initilizing logger: %w", err)
	}

	log, err := logger.InitLogger(&cfg)

	if err != nil {
		logStd.Fatal("error initilizing logger: %w", err)
	}
	defer log.Sync()

	serverAddress := flag.String("a", cfg.ServerAddress, "Адрес сервера")
	baseURL := flag.String("b", cfg.BaseURL, "Базовый адрес")
	fileStoragePath := flag.String("f", cfg.FileStoragePath, "Путь до папки с хранением данных")
	dataBaseDSN := flag.String("d", cfg.DataBaseDSN, "Конфиг подключения к db")

	flag.Parse()

	log.Info("Finished flags env")

	log.Info(
		"Start with params",
		zap.Stringp("serverAddress", serverAddress),
		zap.Stringp("baseURL", baseURL),
		zap.Stringp("fileStoragePath", fileStoragePath),
		zap.Stringp("dataBaseDSN", dataBaseDSN),
	)

	app.NewApp(
		log,
		*serverAddress,
		*baseURL,
		*fileStoragePath,
		*dataBaseDSN,
	)
}

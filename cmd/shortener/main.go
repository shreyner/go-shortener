package main

import (
	"fmt"
	logStd "log"

	"go.uber.org/zap"

	"github.com/shreyner/go-shortener/internal/app"
	"github.com/shreyner/go-shortener/internal/config"
	"github.com/shreyner/go-shortener/internal/pkg/logger"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %v\nBuild date: %v\nBuild commit: %v\n", buildVersion, buildDate, buildCommit)

	var cfg config.Config

	if err := cfg.Parse(); err != nil {
		logStd.Fatal("error initilizing logger: %w", err)
	}

	log, err := logger.InitLogger(&cfg)

	if err != nil {
		logStd.Fatal("error initilizing logger: %w", err)
	}
	defer log.Sync()

	log.Info("Finished flags env")

	log.Info(
		"Start with params",
		zap.String("serverAddress", cfg.ServerAddress),
		zap.String("baseURL", cfg.BaseURL),
		zap.String("fileStoragePath", cfg.FileStoragePath),
		zap.String("dataBaseDSN", cfg.DataBaseDSN),
		zap.Bool("enabledHTTS", cfg.EnabledHTTPS),
		zap.String("config", cfg.Config),
	)

	app.NewApp(log, &cfg)
}

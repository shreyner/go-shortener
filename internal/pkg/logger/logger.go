package logger

import (
	"fmt"
	"github.com/shreyner/go-shortener/internal/config"
	"go.uber.org/zap"
)

func InitLogger(cfg *config.Config) (*zap.Logger, error) {
	cfgLog := zap.NewDevelopmentConfig()

	cfgLog.DisableStacktrace = true

	logger, err := cfgLog.Build()

	if err != nil {
		return nil, fmt.Errorf("logger build: %w", err)
	}

	return logger, nil
}

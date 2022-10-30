// Package logger create and init zap logger
package logger

import (
	"fmt"

	"github.com/shreyner/go-shortener/internal/config"
	"go.uber.org/zap"
)

// InitLogger create and init logger by app config
func InitLogger(cfg *config.Config) (*zap.Logger, error) {
	cfgLog := zap.NewDevelopmentConfig()

	//cfgLog.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)

	cfgLog.DisableStacktrace = true

	logger, err := cfgLog.Build()

	if err != nil {
		return nil, fmt.Errorf("logger build: %w", err)
	}

	return logger, nil
}

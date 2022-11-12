// Package config include all configuration params for application from env variables
//
// Example use:
//
//	var cfg config.Config
//
//	if err := cfg.Parse(); err != nil {
//	  logStd.Fatal("error initilizing logger: %w", err)
//	}
package config

import "github.com/caarlos0/env/v6"

// Config struct has all fields that will be filled in from env vars
type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DataBaseDSN     string `env:"DATABASE_DSN"`
	EnabledHTTPS    bool   `env:"ENABLE_HTTPS"`
}

// Parse will start parsing env variable and willed config
func (c *Config) Parse() error {
	if err := env.Parse(c); err != nil {
		return err
	}

	return nil
}

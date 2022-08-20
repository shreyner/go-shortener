package config

import "github.com/caarlos0/env/v6"

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DataBaseDSN     string `env:"DATABASE_DSN"`
}

func (c *Config) Parse() error {

	if err := env.Parse(c); err != nil {
		return err
	}

	return nil
}

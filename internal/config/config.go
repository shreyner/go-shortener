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

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/caarlos0/env/v6"
)

// Config struct has all fields that will be filled in from env vars
type Config struct {
	ServerAddress   string `json:"server_address" env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `json:"base_url" env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `json:"file_storage_path" env:"FILE_STORAGE_PATH"`
	DataBaseDSN     string `json:"database_dsn" env:"DATABASE_DSN"`
	Config          string `json:"-" env:"CONFIG"`
	TrustedSubnet   string `json:"trusted_subnet"`
	SignKey         string `json:"sign_key" env:"SIGN_KEY" envDefault:"triy6n9rw3"`
	EnabledHTTPS    bool   `json:"enable_https" env:"ENABLE_HTTPS"`
}

// Parse will start parsing env variable and willed config
func (c *Config) Parse() error {
	if err := env.Parse(c); err != nil {
		return err
	}

	flag.StringVar(&c.Config, "c", c.Config, "Путь до файла с конфигурацией")
	flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "Адрес сервера")
	flag.StringVar(&c.BaseURL, "b", c.BaseURL, "Базовый адрес")
	flag.StringVar(&c.FileStoragePath, "f", c.FileStoragePath, "Путь до папки с хранением данных")
	flag.StringVar(&c.DataBaseDSN, "d", c.DataBaseDSN, "Конфиг подключения к db")
	flag.BoolVar(&c.EnabledHTTPS, "s", c.EnabledHTTPS, "HTTPS соединение")
	flag.StringVar(&c.TrustedSubnet, "t", c.TrustedSubnet, "CIDR для доступа к /internal")
	flag.StringVar(&c.SignKey, "sign-key", c.SignKey, "signed cookie key")

	flag.Parse()

	if c.Config != "" {
		return c.ParseConfigFile(c.Config)
	}

	return nil
}

// ParseConfigFile parsed config.json file and merger with this config
func (c *Config) ParseConfigFile(name string) error {
	file, err := os.Open(name)

	if err != nil {
		return err
	}

	data, err := io.ReadAll(file)

	if err != nil {
		return err
	}

	var configJSON Config
	err = json.Unmarshal(data, &configJSON)

	if err != nil {
		return err
	}

	fmt.Printf("config: %#v\n", configJSON)

	if c.DataBaseDSN == "" && configJSON.DataBaseDSN != "" {
		c.DataBaseDSN = configJSON.DataBaseDSN
	}

	if c.ServerAddress == "" && configJSON.ServerAddress != "" {
		c.ServerAddress = configJSON.ServerAddress
	}

	if c.FileStoragePath == "" && configJSON.FileStoragePath != "" {
		c.FileStoragePath = configJSON.FileStoragePath
	}

	if !c.EnabledHTTPS && configJSON.EnabledHTTPS {
		c.EnabledHTTPS = true
	}

	if c.ServerAddress == ":8080" && configJSON.ServerAddress != "" {
		c.ServerAddress = configJSON.ServerAddress
	}

	if c.BaseURL == "http://localhost:8080" && configJSON.BaseURL != "" {
		c.BaseURL = configJSON.BaseURL
	}

	if c.TrustedSubnet == "" && configJSON.TrustedSubnet != "" {
		c.TrustedSubnet = configJSON.TrustedSubnet
	}

	if c.SignKey == "triy6n9rw3" && configJSON.SignKey != "" {
		c.SignKey = configJSON.SignKey
	}

	return nil
}

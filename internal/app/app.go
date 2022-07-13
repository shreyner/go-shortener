package app

import (
	"flag"
	"github.com/shreyner/go-shortener/internal/storage/storage_database"
	"log"

	"github.com/caarlos0/env/v6"

	"github.com/shreyner/go-shortener/internal/handlers"
	"github.com/shreyner/go-shortener/internal/repositories"
	"github.com/shreyner/go-shortener/internal/server"
	"github.com/shreyner/go-shortener/internal/service"
	"github.com/shreyner/go-shortener/internal/storage/storage_file"
	"github.com/shreyner/go-shortener/internal/storage/storage_memory"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DataBaseDSN     string `env:"DATABASE_DSN"`
}

// TODO: Нужно больше логов

func NewApp() {
	log.Println("Start app...")
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)

		return
	}
	log.Println("Finished parse env")

	// TODO: Валидация флагов
	serverAddress := flag.String("a", cfg.ServerAddress, "Адрес сервера")
	baseURL := flag.String("b", cfg.BaseURL, "Базовый адрес")
	fileStoragePath := flag.String("f", cfg.FileStoragePath, "Путь до папки с хранением данных")
	dataBaseDSN := flag.String("d", cfg.DataBaseDSN, "Конфиг подключения к db")

	flag.Parse()
	log.Println("Finished flags env")

	log.Printf(
		"Start with params: serverAddress: %s, baseURL: %s, fileStoragePath: %s\n",
		*serverAddress,
		*baseURL,
		*fileStoragePath,
	)

	// Попробовать спрятать в APP и спрятать за interface.
	// newShortURL repositry
	var shorterRepository repositories.ShortURLRepository

	if *fileStoragePath != "" {
		shorterFileRepository, err := storagefile.NewShortURLStore(*fileStoragePath)

		if err != nil {
			log.Fatal(err)
			return
		}

		defer shorterFileRepository.Close()

		shorterRepository = shorterFileRepository
	} else {
		shorterRepository = storagememory.NewShortURLStore()
	}

	var storageDB storagedatabase.StorageSQL

	if *dataBaseDSN != "" {
		log.Println("Connect to database ...")
		database, err := storagedatabase.NewStorageSQL(*dataBaseDSN)

		if err != nil {
			log.Fatal(err)
			return
		}

		defer database.Close()

		log.Println("Success connected database")

		log.Println("Check and create database...")
		if err := database.CheckAndCreateSchema(); err != nil {
			log.Fatal(err)
			return
		}
		log.Println("Finish check or created...")

		storageDB = database
		shortURLStorage, _ := storagedatabase.NewShortURLStore(database.DB)
		shorterRepository = shortURLStorage
	}

	services := service.NewService(shorterRepository)

	router := handlers.NewRouter(*baseURL, services.ShorterService, storageDB)
	serv := server.NewServer(*serverAddress, router)

	log.Println("Listen server")
	serv.Start()
}

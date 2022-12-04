package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/shreyner/go-shortener/internal/middlewares"
	"github.com/shreyner/go-shortener/internal/rpcservices"
	pb "github.com/shreyner/go-shortener/proto"
	"go.uber.org/zap"

	"github.com/shreyner/go-shortener/internal/config"
	"github.com/shreyner/go-shortener/internal/handlers"
	"github.com/shreyner/go-shortener/internal/pkg/fans"
	"github.com/shreyner/go-shortener/internal/server"
	"github.com/shreyner/go-shortener/internal/service"
	"github.com/shreyner/go-shortener/internal/storage"
)

// NewApp create shortener application and start http listen, db connection and waiting system signal for stop
func NewApp(
	log *zap.Logger,
	cfg *config.Config,
) {
	log.Info("Start app...")

	log.Info("Create storage...")
	store, err := storage.NewStorage(log, cfg.FileStoragePath, cfg.DataBaseDSN)

	if err != nil {
		log.Error("", zap.Error(err))
		return
	}

	//defer func() {
	//	if err := store.Close(); err != nil {
	//		log.Error("error close connection to store", zap.Error(err))
	//	}
	//}()

	log.Info("Create services...")
	services, err := service.NewService(log, store.ShortURL, []byte(cfg.SignKey))

	if err != nil {
		log.Error("can't create services", zap.Error(err))
		return
	}

	log.Info("Create fanShortService...")
	fansShortService := fans.NewFansShortService(log, store.ShortURL, 4)
	//defer fansShortService.Close()

	r := handlers.NewRouter(
		log,
		cfg.BaseURL,
		services.ShorterService,
		services.AuthService,
		store.ShortURL,
		store,
		fansShortService,
		cfg.TrustedSubnet,
	)

	log.Info("Create http server")
	httpserver := server.NewHTTPServer(log, cfg.ServerAddress, r)

	log.Info("Create grpc server")
	grcserver, err := server.NewGRPCServer(log, ":3200", middlewares.AuthInterceptor(services.AuthService))

	if err != nil {
		log.Error("Can't start grpc server", zap.Error(err))
		return
	}

	pb.RegisterAuthServer(grcserver.Server, rpcservices.NewAuthServer(log, services.AuthService))
	pb.RegisterShortenerServer(
		grcserver.Server,
		rpcservices.NewShortenerServer(log, services.ShorterService, fansShortService),
	)

	grcserver.Start()
	httpserver.Start(cfg.EnabledHTTPS)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	select {
	case x := <-interrupt:
		log.Info("Received a signal.", zap.String("signal", x.String()))
	case err := <-httpserver.Notify():
		log.Error("Received an error from the start rest api server", zap.Error(err))
	case err := <-grcserver.Notify():
		log.Error("Received an error from the start grpc server", zap.Error(err))
	}

	log.Info("Stopping server...")

	if err := grcserver.Stop(); err != nil {
		log.Error("Got an error while stopping th grpc server", zap.Error(err))
	}

	if err := httpserver.Stop(context.Background()); err != nil {
		log.Error("Got an error while stopping th rest api server", zap.Error(err))
	}

	fansShortService.Close()

	if err := store.Close(); err != nil {
		log.Error("error close connection to store", zap.Error(err))
	}

	log.Info("The app is calling the last defers and will be stopped.")
}

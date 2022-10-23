package handlers

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/shreyner/go-shortener/internal/pkg/fans"
	"github.com/shreyner/go-shortener/internal/repositories"
	"go.uber.org/zap"

	"github.com/shreyner/go-shortener/internal/middlewares"
	"github.com/shreyner/go-shortener/internal/storage"
)

var cookieSecretKey = []byte("triy6n9rw3")

func NewRouter(
	log *zap.Logger,
	baseURL string,
	shorterService ShortedService,
	shortURIRepository repositories.ShortURLRepository,
	storage *storage.Storage,
	fansShortService *fans.FansShortService,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middlewares.NewStructuredLogger(log))
	r.Use(chiMiddleware.Recoverer)

	authMiddleware := middlewares.AuthHandler(log, cookieSecretKey)

	shortedHandler := NewShortedHandler(log, baseURL, shorterService, shortURIRepository, fansShortService)
	storeHandler := NewStoreHandler(log, storage)

	r.Route("/api", func(r chi.Router) {
		r.With(authMiddleware).Route("/shorten", func(r chi.Router) {
			r.
				With(
					chiMiddleware.AllowContentEncoding("gzip"),
					middlewares.GzlibCompressHandler,
				).
				Post("/", shortedHandler.APICreate)

			r.Post("/batch", shortedHandler.APICreateBatch)
		})

		r.With(authMiddleware).Route("/user", func(r chi.Router) {
			r.Route("/urls", func(r chi.Router) {
				r.Get("/", shortedHandler.APIUserURLs)
				r.Delete("/", shortedHandler.APIUserDeleteURLs)
			})
		})
	})

	r.With(
		chiMiddleware.AllowContentEncoding("gzip"),
		middlewares.GzlibCompressHandler,
		authMiddleware,
	).
		Post("/", shortedHandler.Create)

	r.Get("/{id}", shortedHandler.Get)

	r.Get("/ping", storeHandler.Ping)

	return r

}

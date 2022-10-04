package handlers

import (
	"context"
	"github.com/shreyner/go-shortener/internal/pkg/fans"
	"github.com/shreyner/go-shortener/internal/repositories"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
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
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	shortedHandler := NewShortedHandler(log, baseURL, shorterService, shortURIRepository, fansShortService)

	r.Route("/api", func(r chi.Router) {
		r.With(middlewares.AuthHandler(cookieSecretKey)).Route("/shorten", func(r chi.Router) {
			r.
				With(
					chiMiddleware.AllowContentEncoding("gzip"),
					middlewares.GzlibCompressHandler,
				).
				Post("/", shortedHandler.APICreate)

			r.Post("/batch", shortedHandler.APICreateBatch)
		})

		r.With(middlewares.AuthHandler(cookieSecretKey)).Route("/user", func(r chi.Router) {
			r.Route("/urls", func(r chi.Router) {
				r.Get("/", shortedHandler.APIUserURLs)
				r.Delete("/", shortedHandler.APIUserDeleteURLs)
			})
		})
	})

	r.With(
		chiMiddleware.AllowContentEncoding("gzip"),
		middlewares.GzlibCompressHandler,
		middlewares.AuthHandler(cookieSecretKey),
	).
		Post("/", shortedHandler.Create)

	r.Get("/{id}", shortedHandler.Get)

	r.Get("/ping", func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := storage.PingContext(ctx); err != nil {
			log.Error("can't ping to database", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)

			return
		}

		rw.WriteHeader(http.StatusOK)
	})

	return r

}

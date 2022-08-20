package router

import (
	"context"
	"github.com/shreyner/go-shortener/internal/storage"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/shreyner/go-shortener/internal/handlers"
	"github.com/shreyner/go-shortener/internal/middlewares"
)

var cookieSecretKey = []byte("triy6n9rw3")

func NewRouter(
	log *zap.Logger,
	baseURL string,
	shorterService handlers.ShortedService,

	storage *storage.Storage,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	shortedHandler := handlers.NewShortedHandler(baseURL, shorterService)

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
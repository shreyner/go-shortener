package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/shreyner/go-shortener/internal/middlewares"
	storagedatabase "github.com/shreyner/go-shortener/internal/storage/storage_database"
)

var cookieSecretKey = []byte("triy6n9rw3")

func NewRouter(baseURL string, shorterService ShortedService, storageDB storagedatabase.StorageSQL) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	// handlers
	shortedHandler := NewShortedHandler(baseURL, shorterService)

	r.Route("/api", func(r chi.Router) {
		r.
			With(
				chiMiddleware.AllowContentEncoding("gzip"),
				middlewares.GzlibCompressHandler,
				middlewares.AuthHandler(cookieSecretKey),
			).
			Post("/shorten", shortedHandler.APICreate)

		r.Route("/user", func(r chi.Router) {
			r.With(middlewares.AuthHandler(cookieSecretKey)).Get("/urls", shortedHandler.APIUserURLs)
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

		if err := storageDB.PingContext(ctx); err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
	})

	return r
}

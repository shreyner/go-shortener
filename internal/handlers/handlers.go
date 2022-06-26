package handlers

import (
	"github.com/shreyner/go-shortener/internal/middlewares"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

var cookieSecretKey = []byte("triy6n9rw3")

func NewRouter(baseURL string, shorterService ShortedService) *chi.Mux {
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

	return r
}

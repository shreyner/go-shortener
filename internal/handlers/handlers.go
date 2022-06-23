package handlers

import (
	"github.com/shreyner/go-shortener/internal/middlewares"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

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
			With(chiMiddleware.AllowContentEncoding("gzip"), middlewares.GzlibCompressHandler).
			Post("/shorten", shortedHandler.APICreate)
	})
	r.With(chiMiddleware.AllowContentEncoding("gzip"), middlewares.GzlibCompressHandler).
		Post("/", shortedHandler.Create)
	r.Get("/{id}", shortedHandler.Get)

	return r
}

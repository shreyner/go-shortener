package handlers

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	middleware "github.com/shreyner/go-shortener/internal/middlewares"
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
			With(chiMiddleware.AllowContentEncoding("gzip"), middleware.GzipCompressHandler).
			Post("/shorten", shortedHandler.APICreate)
	})
	r.With(chiMiddleware.AllowContentEncoding("gzip"), middleware.GzipCompressHandler).
		Post("/", shortedHandler.Create)
	r.Get("/{id}", shortedHandler.Get)

	return r
}

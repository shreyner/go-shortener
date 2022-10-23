package handlers

import (
	"compress/gzip"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/shreyner/go-shortener/internal/pkg/fans"
	"github.com/shreyner/go-shortener/internal/repositories"
	"go.uber.org/zap"

	"github.com/shreyner/go-shortener/internal/middlewares"
)

var cookieSecretKey = []byte("triy6n9rw3")

func NewRouter(
	log *zap.Logger,
	baseURL string,
	shorterService ShortedService,
	shortURIRepository repositories.ShortURLRepository,
	fansShortService *fans.FansShortService,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middlewares.NewStructuredLogger(log))
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Compress(gzip.BestSpeed, "gzip"))

	authMiddleware := middlewares.AuthHandler(log, cookieSecretKey)

	shortedHandler := NewShortedHandler(log, baseURL, shorterService, shortURIRepository, fansShortService)

	r.Route("/api", func(r chi.Router) {
		r.With(
			chiMiddleware.AllowContentType("application/json"),
			authMiddleware,
		).
			Group(func(r chi.Router) {
				r.Route("/shorten", func(r chi.Router) {
					r.Post("/", shortedHandler.APICreate)
					r.Post("/batch", shortedHandler.APICreateBatch)
				})

				r.Route("/user", func(r chi.Router) {
					r.Route("/urls", func(r chi.Router) {
						r.Get("/", shortedHandler.APIUserURLs)
						r.Delete("/", shortedHandler.APIUserDeleteURLs)
					})
				})
			})
	})

	r.With(
		chiMiddleware.AllowContentType("text/plain"),
		authMiddleware,
	).
		Post("/", shortedHandler.Create)

	r.Get("/{id}", shortedHandler.Get)

	return r
}

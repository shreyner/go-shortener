package handlers

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/shreyner/go-shortener/internal/middlewares"
	"github.com/shreyner/go-shortener/internal/pkg/fans"
	"github.com/shreyner/go-shortener/internal/repositories"
	"github.com/shreyner/go-shortener/internal/storage"
)

// @title       Shortener API
// @description Сервис сокращения ссылок
// @version     1.0

// @host localhost:8080

type authService interface {
	GenerateUserID() string
	CreateToken(userID string) string
	GetUserIDFromToken(token string) (string, error)
}

// NewRouter init and create all handler on route
func NewRouter(
	log *zap.Logger,
	baseURL string,
	shorterService ShortedService,
	authService authService,
	shortURIRepository repositories.ShortURLRepository,
	storage *storage.Storage,
	fansShortService *fans.FansShortService,
	trustedSubnet string,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middlewares.NewStructuredLogger(log))
	r.Use(chiMiddleware.Recoverer)

	authMiddleware := middlewares.AuthHandler(authService)
	realIPMiddleware := middlewares.RealIP
	cidrAccessMiddleware, _ := middlewares.CIDRAccess(trustedSubnet) // 192.168.88.0/24,127.0.0.1/32

	shortedHandler := NewShortedHandler(log, baseURL, shorterService, shortURIRepository, fansShortService)
	storeHandler := NewStoreHandler(log, storage)
	internalHandler := NewInternalHandler(log, shortURIRepository)

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

		r.With(realIPMiddleware, cidrAccessMiddleware).Route("/internal", func(r chi.Router) {
			r.Get("/stats", internalHandler.GetStats)
		})
	})

	r.With(
		chiMiddleware.AllowContentEncoding("gzip"),
		middlewares.GzlibCompressHandler,
		authMiddleware,
	).
		Post("/", shortedHandler.Create)

	r.Get("/ping", storeHandler.Ping)

	r.Get("/{id}", shortedHandler.Get)

	//r.Mount("/debug", chiMiddleware.Profiler())

	return r
}

package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(shorterService ShortedService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// handlers
	shortedHandler := NewShortedHandler(shorterService)

	r.Post("/", shortedHandler.Create)
	r.Get("/{id}", shortedHandler.Get)

	return r
}

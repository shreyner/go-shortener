package handlers

import (
	"context"
	"github.com/shreyner/go-shortener/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type StoreHandler struct {
	log   *zap.Logger
	store *storage.Storage
}

func NewStoreHandler(log *zap.Logger, store *storage.Storage) *StoreHandler {
	return &StoreHandler{
		log:   log,
		store: store,
	}
}

func (s *StoreHandler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := s.store.PingContext(ctx); err != nil {
		s.log.Error("can't ping to database", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

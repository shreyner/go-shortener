package handlers

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Storage interface {
	PingContext(ctx context.Context) error
}

type StoreHandler struct {
	log   *zap.Logger
	store Storage
}

func NewStoreHandler(log *zap.Logger, store Storage) *StoreHandler {
	return &StoreHandler{
		log:   log,
		store: store,
	}
}

// Ping Состояние подключения к DB
//
// @summary     Состояние подключения к DB
// @description Проверка состояние подключения к базе данных
// @tags        HealthAPI
// @success     200
// @failure     500
// @router      /ping [get]
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
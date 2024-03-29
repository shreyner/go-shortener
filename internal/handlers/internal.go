package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/shreyner/go-shortener/internal/core"
	"go.uber.org/zap"
)

type internalRepository interface {
	GetStats(ctx context.Context) (*core.ShortStats, error)
}

// InternalHandler with internal handlers
type InternalHandler struct {
	log        *zap.Logger
	repository internalRepository
}

// NewInternalHandler create struct InternalHandler
func NewInternalHandler(log *zap.Logger, repository internalRepository) *InternalHandler {
	return &InternalHandler{
		log:        log,
		repository: repository,
	}
}

// GetStats handler stats endpoint
func (i *InternalHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stats, err := i.repository.GetStats(ctx)

	if err != nil {
		i.log.Error("get stats error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	body, err := json.Marshal(stats)

	if err != nil {
		i.log.Error("json marshal stats error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

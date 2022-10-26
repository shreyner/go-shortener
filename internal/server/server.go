package server

import (
	"context"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	server http.Server
	errors chan error
	log    *zap.Logger
}

func NewServer(log *zap.Logger, address string, router http.Handler) *Server {
	return &Server{
		server: http.Server{
			Addr:    address,
			Handler: router,
		},
		log:    log,
		errors: make(chan error),
	}
}

func (s *Server) Start() {
	go func() {
		s.log.Info("Http Server listening on ", zap.String("addr", s.server.Addr))
		s.errors <- s.server.ListenAndServe()
		close(s.errors)
	}()
}

func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return s.server.Shutdown(ctx)
}

func (s *Server) Notify() <-chan error {
	return s.errors
}

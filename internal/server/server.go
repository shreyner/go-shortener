// Package server create http transport
//
// For example create:
//
//	serv := server.NewServer(log, serverAddress, r)
//	serv.Start()
//	serv.Stop()
package server

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

// Server http server
type Server struct {
	server http.Server
	errors chan error
	log    *zap.Logger
}

// NewServer create http transport
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

// Start http server and listening port
func (s *Server) Start() {
	go func() {
		s.log.Info("Http Server listening on ", zap.String("addr", s.server.Addr))
		s.errors <- s.server.ListenAndServe()
		close(s.errors)
	}()
}

// Stop don't listen new connection and waiting close current active connection. Used for graceful shutdown
func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return s.server.Shutdown(ctx)
}

// Notify return chan with error. For example include error "port already exists"
func (s *Server) Notify() <-chan error {
	return s.errors
}

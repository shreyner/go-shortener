package server

import (
	"log"
	"net/http"
)

type Server struct {
	server http.Server
}

func NewServer(address string, router http.Handler) *Server {
	return &Server{
		server: http.Server{
			Addr:    address,
			Handler: router,
		},
	}
}

func (s *Server) Start() {
	log.Fatalln(s.server.ListenAndServe())
}

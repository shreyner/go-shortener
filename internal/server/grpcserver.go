package server

import (
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	log    *zap.Logger
	errors chan error

	Server *grpc.Server
	listen net.Listener
}

func NewGRPCServer(log *zap.Logger, address string, interceptors ...grpc.UnaryServerInterceptor) (*GRPCServer, error) {
	grpcServer := GRPCServer{
		log:    log,
		errors: make(chan error),

		Server: grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors...)),
	}

	listen, err := net.Listen("tcp", address)

	if err != nil {
		return nil, err
	}

	grpcServer.listen = listen

	return &grpcServer, nil
}

func (s *GRPCServer) Start() {
	go func() {
		s.log.Info("gRPC server listen on ", zap.String("addr", s.listen.Addr().String()))
		defer close(s.errors)

		s.errors <- s.Server.Serve(s.listen)
	}()
}

func (s *GRPCServer) Stop() error {
	s.log.Info("gRPC server stopped...")
	s.Server.GracefulStop()

	return s.listen.Close()
}

func (s *GRPCServer) Notify() <-chan error {
	return s.errors
}

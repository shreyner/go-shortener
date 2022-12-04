package rpcservices

import (
	"context"

	"go.uber.org/zap"

	"github.com/shreyner/go-shortener/internal/service"
	pb "github.com/shreyner/go-shortener/proto"
)

var (
	_ pb.AuthServer = (*AuthServer)(nil)
)

// AuthServer base auth grpc server
type AuthServer struct {
	pb.UnimplementedAuthServer

	log         *zap.Logger
	authService *service.AuthService
}

// NewAuthServer constructor
func NewAuthServer(log *zap.Logger, authService *service.AuthService) *AuthServer {
	return &AuthServer{
		log:         log,
		authService: authService,
	}
}

// GetToken create token for identify next requests
func (s *AuthServer) GetToken(_ context.Context, _ *pb.Empty) (*pb.GetTokenResponse, error) {
	var getTokenResponse pb.GetTokenResponse

	getTokenResponse.Token = s.authService.CreateToken(s.authService.GenerateUserID())

	return &getTokenResponse, nil
}

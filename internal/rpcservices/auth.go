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

type AuthServer struct {
	pb.UnimplementedAuthServer

	log         *zap.Logger
	authService *service.AuthService
}

func NewAuthServer(log *zap.Logger, authService *service.AuthService) *AuthServer {
	return &AuthServer{
		log:         log,
		authService: authService,
	}
}

func (s *AuthServer) GetToken(_ context.Context, _ *pb.Empty) (*pb.GetTokenResponse, error) {
	var getTokenResponse pb.GetTokenResponse

	getTokenResponse.Token = s.authService.CreateToken(s.authService.GenerateUserID())

	return &getTokenResponse, nil
}

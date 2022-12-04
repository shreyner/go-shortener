package rpcservices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/shreyner/go-shortener/internal/core"
	"github.com/shreyner/go-shortener/internal/middlewares"
	"github.com/shreyner/go-shortener/internal/pkg/fans"
	sdb "github.com/shreyner/go-shortener/internal/storage/store_errors"
	pb "github.com/shreyner/go-shortener/proto"
)

var (
	_ pb.ShortenerServer = (*ShortenerServer)(nil)
)

type shortedService interface {
	Create(ctx context.Context, userID, url string) (*core.ShortURL, error)
	CreateBatch(ctx context.Context, shortURLs *[]*core.ShortURL) error
	GetByID(ctx context.Context, key string) (*core.ShortURL, bool)
	AllByUser(ctx context.Context, id string) ([]*core.ShortURL, error)
}

// ShortenerServer base shortner handler for grpc server
type ShortenerServer struct {
	pb.UnimplementedShortenerServer

	log              *zap.Logger
	service          shortedService
	fansShortService *fans.FansShortService
}

// NewShortenerServer constructor
func NewShortenerServer(
	log *zap.Logger,
	service shortedService,
	fansShortService *fans.FansShortService,
) *ShortenerServer {
	return &ShortenerServer{
		log:              log,
		service:          service,
		fansShortService: fansShortService,
	}
}

// CreateShort create new short by url
func (s *ShortenerServer) CreateShort(
	ctx context.Context,
	in *pb.CreateShortRequest,
) (*pb.CreateShortResponse, error) {
	userID, _ := middlewares.GetUserIDCtx(ctx)
	var response pb.CreateShortResponse

	if in.Url == "" {
		response.Error = "url is required"
		return &response, nil
	}

	if _, err := url.Parse(in.Url); err != nil {
		response.Error = "invalid url"
		return &response, nil
	}

	shortURL, err := s.service.Create(ctx, userID, in.Url)

	var shortURLCreateConflictError *sdb.ShortURLCreateConflictError
	if errors.As(err, &shortURLCreateConflictError) {
		response.Id = shortURLCreateConflictError.OriginID

		return &response, nil
	}

	if err != nil {
		s.log.Error("unhandled error when create short url", zap.Error(err))
		return nil, status.Error(codes.Internal, "unhandled error when create short url")
	}

	response.Id = shortURL.ID

	return &response, nil
}

// CreateBatchShort create batch sorts by urls
func (s *ShortenerServer) CreateBatchShort(
	ctx context.Context,
	in *pb.CreateBatchShortRequest,
) (*pb.CreateBatchShortResponse, error) {
	userID, _ := middlewares.GetUserIDCtx(ctx)
	var response pb.CreateBatchShortResponse

	shoredURLs := make([]*core.ShortURL, len(in.Urls))

	for i, v := range in.Urls {
		_, err := url.Parse(v.Url)

		if err != nil {
			response.Error = fmt.Sprintf("invalid url for CorrelationID: %v", v.CorrelationId)
			return &response, nil
		}

		shortURL := core.ShortURL{
			URL: v.Url,
			UserID: sql.NullString{
				String: userID,
				Valid:  userID != "",
			},
			CorrelationID: v.CorrelationId,
		}

		shoredURLs[i] = &shortURL
	}

	if err := s.service.CreateBatch(ctx, &shoredURLs); err != nil {
		s.log.Error("unhandled error when create short url", zap.Error(err))
		return nil, status.Error(codes.Internal, "unhandled error when create short url")
	}

	responseURLs := make([]*pb.CreateBatchShortResponse_URL, len(shoredURLs))

	for i, v := range shoredURLs {
		responseURL := pb.CreateBatchShortResponse_URL{
			Id:            v.ID,
			CorrelationId: v.CorrelationID,
		}

		responseURLs[i] = &responseURL
	}

	response.Urls = responseURLs

	return &response, nil
}

// ListUserURLs return list shorted was created user
func (s *ShortenerServer) ListUserURLs(
	ctx context.Context,
	_ *pb.ListUserURLsRequest,
) (*pb.ListUserURLsResponse, error) {
	userID, ok := middlewares.GetUserIDCtx(ctx)
	var listUserURLsResponse pb.ListUserURLsResponse

	if !ok || userID == "" {
		s.log.Info("missing token", zap.String("userID", userID))
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}

	list, err := s.service.AllByUser(ctx, userID)

	if err != nil {
		s.log.Error("unhandled error when get list urls by userID", zap.Error(err))
		return nil, status.Error(codes.Internal, "unhandled error")
	}

	responseList := make([]*pb.ListUserURLsResponse_URL, len(list))

	for i, shortURL := range list {
		responseURL := pb.ListUserURLsResponse_URL{
			Id:          shortURL.ID,
			OriginalURL: shortURL.URL,
		}

		responseList[i] = &responseURL
	}

	listUserURLsResponse.Urls = responseList

	return &listUserURLsResponse, nil
}

// DeleteByIDs delete by ids for current user
func (s *ShortenerServer) DeleteByIDs(
	ctx context.Context,
	in *pb.DeleteByIDsRequest,
) (*pb.DeleteByIDsResponse, error) {
	userID, ok := middlewares.GetUserIDCtx(ctx)
	var deleteByIDsResponse pb.DeleteByIDsResponse

	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}

	s.log.Info("was delete", zap.String("userID", userID), zap.Strings("urlIDs", in.Ids))

	s.fansShortService.Add(userID, in.Ids)

	return &deleteByIDsResponse, nil
}

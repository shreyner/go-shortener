package middlewares

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/shreyner/go-shortener/internal/service"
)

// UserCtxKey uniq key for save user ID in context
type UserCtxKey int

const userCtxKey UserCtxKey = iota

var (
	authCookieKey = "auth"
	tokenKey      = "token"
)

type authService interface {
	GenerateUserID() string
	CreateToken(userID string) string
	GetUserIDFromToken(token string) (string, error)
}

// GetUserIDCtx return user ID from context
func GetUserIDCtx(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(userCtxKey).(string)
	return v, ok
}

func SetUserIDCtx(parentCtx context.Context, userID string) context.Context {
	return context.WithValue(parentCtx, userCtxKey, userID)
}

// AuthHandler for auth users and create if not found auth cookies
func AuthHandler(authService authService) func(next http.Handler) http.Handler {
	parseCookie := func(r *http.Request) (string, error) {
		authCookie, err := r.Cookie(authCookieKey)

		if err != nil {
			return "", err
		}

		userID, err := authService.GetUserIDFromToken(authCookie.Value)

		if err != nil {
			return "", err
		}

		return userID, nil
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			userID, err := parseCookie(r)
			ctx := r.Context()

			if err != nil {
				newUserID := authService.GenerateUserID()

				dst := authService.CreateToken(newUserID)

				authCookie := &http.Cookie{Value: dst, Name: authCookieKey}

				http.SetCookie(rw, authCookie)

				ctx = SetUserIDCtx(ctx, newUserID)

				next.ServeHTTP(rw, r.WithContext(ctx))

				return
			}

			ctx = SetUserIDCtx(ctx, userID)

			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

func AuthInterceptor(authService *service.AuthService) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)

		if !ok {
			return handler(ctx, req)
		}

		var token string

		if value := md.Get(tokenKey); len(value) > 0 {
			token = value[0]
		}

		if token == "" {
			return handler(ctx, req)
		}

		userID, err := authService.GetUserIDFromToken(token)

		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		return handler(SetUserIDCtx(ctx, userID), req)
	}
}

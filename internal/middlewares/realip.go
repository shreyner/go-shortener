package middlewares

import (
	"context"
	"net"
	"net/http"
	"strings"
)

// RealIPCtxKey type for key real ip in context
type RealIPCtxKey int

const realIPCtxKey RealIPCtxKey = iota

var (
	xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
	xReadIP       = http.CanonicalHeaderKey("X-Read-IP")
)

func GetRealIPCtx(ctx context.Context) net.IP {
	v, _ := ctx.Value(realIPCtxKey).(net.IP)

	return v
}

func setReadIPCtx(ctx context.Context, ip net.IP) context.Context {
	return context.WithValue(ctx, realIPCtxKey, ip)
}

func RealIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		reqIPAddr := getRealIPRequest(r)
		realIPAddr := net.ParseIP(reqIPAddr)

		// TODO: Попробовать сделать на ссылку
		next.ServeHTTP(w, r.WithContext(setReadIPCtx(ctx, realIPAddr)))
	})
}

func CIDRAccess(CIDRs string) (func(http.Handler) http.Handler, error) {
	var ipNet *net.IPNet

	if CIDRs != "" {
		_, ipNetParsed, err := net.ParseCIDR(CIDRs)

		if err != nil {
			return nil, err
		}

		ipNet = ipNetParsed
	}

	handler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			realIP := GetRealIPCtx(r.Context())

			if realIP == nil || ipNet == nil {
				// TODO: Add logs

				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)

				return
			}

			if !ipNet.Contains(realIP) {
				// TODO: Add logs

				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)

				return
			}

			next.ServeHTTP(w, r)
		})
	}

	return handler, nil
}

func getRealIPRequest(r *http.Request) string {
	if xrip := r.Header.Get(xReadIP); xrip != "" {
		return xrip
	}

	if xff := r.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ", ")

		if i == -1 {
			i = len(xff)
		}

		return xff[:i]
	}

	return ""
}

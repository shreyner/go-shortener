package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type StructuredLogger struct {
	Logger *zap.Logger
}

func (s *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: s.Logger.Named("Http Request")}
	fields := make([]zap.Field, 0)

	//fields = append(fields, zap.String("ts", time.Now().UTC().Format(time.RFC1123)))

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		fields = append(fields, zap.String("req_id", reqID))
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	fields = append(fields, zap.String("http_scheme", scheme))
	fields = append(fields, zap.String("http_method", r.Method))
	fields = append(fields, zap.String("remote_addr", r.RemoteAddr))
	fields = append(fields, zap.String("user_agent", r.UserAgent()))
	fields = append(fields, zap.String("uri", fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)))

	entry.Logger = entry.Logger.With(fields...)

	entry.Logger.Info("request started")

	return entry
}

type StructuredLoggerEntry struct {
	Logger *zap.Logger
}

func (s *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	log := s.Logger.With(
		zap.Int("resp_status", status),
		zap.Int("resp_bytes_length", bytes),
		zap.Float64("resp_elapsed_ms", float64(elapsed.Nanoseconds())/1000000.0), //nolint:gomnd
	)

	switch {
	case status <= http.StatusBadRequest:
		log.Info("request complete")
	case status <= http.StatusInternalServerError:
		log.Warn("request complete")
	default:
		log.Error("request complete")
	}
}

func (s *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	s.Logger = s.Logger.With(
		zap.ByteString("stack", stack),
		zap.String("panic", fmt.Sprintf("%+v", v)),
	)
}

func NewStructuredLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

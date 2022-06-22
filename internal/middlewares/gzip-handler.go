package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipCompressHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)

			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		defer gz.Close()

		w.Header().Add("Content-Encoding", "gzip")

		next.ServeHTTP(gzipWriter{w, gz}, r)
	})
}
package middlewares

import (
	"compress/zlib"
	"io"
	"net/http"
	"strings"
)

type gzlibWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzlibWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzlibCompressHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)

			return
		}

		zw, err := zlib.NewWriterLevel(w, zlib.BestSpeed)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		defer zw.Close()

		w.Header().Add("Content-Encoding", "gzip")

		next.ServeHTTP(gzlibWriter{w, zw}, r)
	})
}

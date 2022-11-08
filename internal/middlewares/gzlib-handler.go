package middlewares

import (
	"compress/gzip"
	"io"
	"log"
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

// GzlibCompressHandler compress response
func GzlibCompressHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(rw, r)

			return
		}

		zw, err := gzip.NewWriterLevel(rw, gzip.BestSpeed)

		if err != nil {
			log.Println("error compress response", err.Error())
			http.Error(rw, err.Error(), http.StatusInternalServerError)

			return
		}

		defer zw.Close()

		rw.Header().Add("Content-Encoding", "gzip")

		// TODO: Вынести сюда логику с decompress, разбирать body и складывать stream buffer в body

		next.ServeHTTP(gzlibWriter{rw, zw}, r)
	})
}

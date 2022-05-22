package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

var listenServerAddress = ":8080"

var mapShortedUrls = map[string]string{
	"g": "https://google.com",
	"y": "https://yandex.ru",
}

func main() {
	mux := http.NewServeMux()

	// TODO: Обработать такие /ыфв/фыв
	// TODO: Обработать не валидный
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		urlPart := removeEmptyStrings(strings.Split(r.URL.Path, "/"))

		if len(urlPart) == 0 && r.Method != http.MethodPost {
			http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		if len(urlPart) == 1 && r.Method != http.MethodGet {
			http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		if len(urlPart) > 1 {
			http.NotFound(w, r)
			return
		}

		if r.Method == http.MethodPost {
			mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))

			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			if mediaType != "application/x-www-form-urlencoded" {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			body, err := io.ReadAll(r.Body)

			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			_, err = url.ParseRequestURI(string(body))
			log.Println("err", err)

			if err != nil {
				http.Error(w, "Invalid url", http.StatusBadRequest)
				return
			}

			key := randSeq(4)
			mapShortedUrls[key] = string(body)

			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, key)

			return
		}

		if r.Method == http.MethodGet {
			shortCode := urlPart[0]

			url, ok := mapShortedUrls[shortCode]

			if !ok {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}

			http.Redirect(w, r, url, http.StatusPermanentRedirect)
			return
		}

		fmt.Fprintln(w, "Hello world")
	})

	log.Printf("Start start on %s", listenServerAddress)
	log.Fatalln(http.ListenAndServe(listenServerAddress, mux))
}

func removeEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

package handlers

import (
	"fmt"
	"github.com/shreyner/go-shortener/internal/core"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

type ShortedService interface {
	Create(url string) (core.ShortURL, error)
	GetByID(key string) (core.ShortURL, bool)
}

type ShortedHandler struct {
	ShorterService ShortedService
}

func NewShortedHandler(shorterService ShortedService) *ShortedHandler {
	return &ShortedHandler{ShorterService: shorterService}
}

func (sh *ShortedHandler) ShortedCreate(wr http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))

	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	if mediaType != "text/plain" {
		http.Error(wr, "bad request", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = url.ParseRequestURI(string(body))

	if err != nil {
		http.Error(wr, "Invalid url", http.StatusBadRequest)
		return
	}

	shortURL, err := sh.ShorterService.Create(string(body))
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	wr.WriteHeader(http.StatusCreated)
	fmt.Fprintf(wr, "http://localhost:8080/%s", shortURL.ID)
}

func (sh *ShortedHandler) ShortedGet(wr http.ResponseWriter, r *http.Request) {
	urlPart := removeEmptyStrings(strings.Split(r.URL.Path, "/"))
	shortCode := urlPart[0]

	shortURL, ok := sh.ShorterService.GetByID(shortCode)

	if !ok {
		http.Error(wr, "Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(wr, r, shortURL.URL, http.StatusTemporaryRedirect)
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

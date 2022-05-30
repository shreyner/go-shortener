package handlers

import (
	"github.com/shreyner/go-shortener/internal/core"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

type ShortedService interface {
	Create(url string) (core.ShortUrl, error)
	GetById(key string) (core.ShortUrl, bool)
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

	shortUrl, err := sh.ShorterService.Create(string(body))
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	wr.WriteHeader(http.StatusCreated)
	wr.Write([]byte(shortUrl.Id))
}

func (sh *ShortedHandler) ShortedGet(wr http.ResponseWriter, r *http.Request) {
	urlPart := removeEmptyStrings(strings.Split(r.URL.Path, "/"))
	shortCode := urlPart[0]

	shortedUrl, ok := sh.ShorterService.GetById(shortCode)

	if !ok {
		http.Error(wr, "Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(wr, r, shortedUrl.Url, http.StatusPermanentRedirect)
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

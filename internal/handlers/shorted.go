package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"

	core "github.com/shreyner/go-shortener/internal/core"
	"github.com/timewasted/go-accept-headers"

	"github.com/go-chi/chi/v5"
)

var (
	CONTENT_TYPE_JSON = "application/json"
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

func (sh *ShortedHandler) Create(wr http.ResponseWriter, r *http.Request) {
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

func (sh *ShortedHandler) Get(wr http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "id")

	shortURL, ok := sh.ShorterService.GetByID(shortCode)

	if !ok {
		http.Error(wr, "Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(wr, r, shortURL.URL, http.StatusTemporaryRedirect)
}

type ShortedCreateDTO struct {
	Url string `json:"url"`
}

type ShortedResponseDTO struct {
	Result string `json:"result"`
}

func (sh *ShortedHandler) ApiCreate(wr http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))

	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	if mediaType != CONTENT_TYPE_JSON {
		http.Error(wr, "bad request", http.StatusBadRequest)
		return
	}

	acceptHeader := r.Header.Get("Accept")

	if acceptHeader != "" {
		crossAccepting, err := accept.Negotiate(acceptHeader, CONTENT_TYPE_JSON)

		if err != nil {
			http.Error(wr, "bad headers", http.StatusBadRequest)
			return
		}

		if crossAccepting != CONTENT_TYPE_JSON {
			http.Error(wr, "bad accepting content", http.StatusNotAcceptable)
			return
		}
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	var shortedCreateDTO ShortedCreateDTO

	if err := json.Unmarshal(body, &shortedCreateDTO); err != nil {
		http.Error(wr, "Error parse body", http.StatusInternalServerError)
		return
	}

	if _, err := url.ParseRequestURI(string(shortedCreateDTO.Url)); err != nil {
		http.Error(wr, "Invalid url", http.StatusBadRequest)
		return
	}

	shortURL, err := sh.ShorterService.Create(shortedCreateDTO.Url)

	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	resulturl := fmt.Sprintf("http://localhost:8080/%s", shortURL.ID)

	responseCreateDTO := ShortedResponseDTO{Result: resulturl}

	responseBody, err := json.Marshal(responseCreateDTO)

	if err != nil {
		http.Error(wr, "error create response", http.StatusInternalServerError)
		return
	}

	wr.WriteHeader(http.StatusCreated)
	wr.Header().Add("Content-type", CONTENT_TYPE_JSON)
	fmt.Fprint(wr, string(responseBody))
}

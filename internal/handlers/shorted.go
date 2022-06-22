package handlers

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"

	core "github.com/shreyner/go-shortener/internal/core"
	"github.com/timewasted/go-accept-headers"

	"github.com/go-chi/chi/v5"
)

var (
	ContentTypeJSON = "application/json"
)

type ShortedService interface {
	Create(url string) (*core.ShortURL, error)
	GetByID(key string) (*core.ShortURL, bool)
}

type ShortedHandler struct {
	ShorterService ShortedService
	baseURL        string
}

func NewShortedHandler(baseURL string, shorterService ShortedService) *ShortedHandler {
	return &ShortedHandler{ShorterService: shorterService, baseURL: baseURL}
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
	fmt.Fprintf(wr, "%s/%s", sh.baseURL, shortURL.ID)
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
	URL string `json:"url"`
}

type ShortedResponseDTO struct {
	Result string `json:"result"`
}

func (sh *ShortedHandler) APICreate(wr http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))

	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	if mediaType != ContentTypeJSON {
		http.Error(wr, "bad request", http.StatusBadRequest)
		return
	}

	acceptHeader := r.Header.Get("Accept")

	if acceptHeader != "" {
		crossAccepting, err := accept.Negotiate(acceptHeader, ContentTypeJSON)

		if err != nil {
			http.Error(wr, "bad headers", http.StatusBadRequest)
			return
		}

		if crossAccepting != ContentTypeJSON {
			http.Error(wr, "bad accepting content", http.StatusNotAcceptable)
			return
		}
	}

	var body []byte

	if strings.Contains(r.Header.Get("Content-Type"), "gzip") {
		if body, err = Decompress(r.Body); err != nil {
			http.Error(wr, err.Error(), http.StatusInternalServerError)

			return
		}
	} else {
		body, err = io.ReadAll(r.Body)
		if err != nil {
			http.Error(wr, err.Error(), http.StatusInternalServerError)

			return
		}
	}
	defer r.Body.Close()

	var shortedCreateDTO ShortedCreateDTO

	if err := json.Unmarshal(body, &shortedCreateDTO); err != nil {
		http.Error(wr, "Error parse body", http.StatusInternalServerError)
		return
	}

	if _, err := url.ParseRequestURI(shortedCreateDTO.URL); err != nil {
		http.Error(wr, "Invalid url", http.StatusBadRequest)
		return
	}

	shortURL, err := sh.ShorterService.Create(shortedCreateDTO.URL)

	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	resultURL := fmt.Sprintf("%s/%s", sh.baseURL, shortURL.ID)

	responseCreateDTO := ShortedResponseDTO{Result: resultURL}

	responseBody, err := json.Marshal(responseCreateDTO)

	if err != nil {
		http.Error(wr, "error create response", http.StatusInternalServerError)
		return
	}

	wr.Header().Add("Content-Type", "application/json")
	wr.WriteHeader(http.StatusCreated)

	wr.Write(responseBody)
}

func Decompress(dateRead io.Reader) ([]byte, error) {
	r := flate.NewReader(dateRead)
	defer r.Close()

	var b bytes.Buffer

	if _, err := b.ReadFrom(r); err != nil {
		return nil, fmt.Errorf("failed decopress data :%w", err)
	}

	return b.Bytes(), nil
}

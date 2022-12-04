package handlers

import (
	"bytes"
	"compress/gzip"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/shreyner/go-shortener/internal/core"
	"github.com/shreyner/go-shortener/internal/middlewares"
	"github.com/shreyner/go-shortener/internal/pkg/fans"
	"github.com/shreyner/go-shortener/internal/pkg/pool"
	"github.com/shreyner/go-shortener/internal/repositories"
	sdb "github.com/shreyner/go-shortener/internal/storage/store_errors"
	"github.com/timewasted/go-accept-headers"
	"go.uber.org/zap"
)

var (
	contentTypeJSON = "application/json"
)

// ShortedService interface for service with business logic
type ShortedService interface {
	Create(ctx context.Context, userID, url string) (*core.ShortURL, error)
	CreateBatch(ctx context.Context, shortURLs *[]*core.ShortURL) error
	GetByID(ctx context.Context, key string) (*core.ShortURL, bool)
	AllByUser(ctx context.Context, id string) ([]*core.ShortURL, error)
}

// ShortedHandler include handlers for shorteners handlers
type ShortedHandler struct {
	log               *zap.Logger
	ShorterService    ShortedService
	ShorterRepository repositories.ShortURLRepository
	fansShortService  *fans.FansShortService
	baseURL           string
}

// NewShortedHandler create instance
func NewShortedHandler(
	log *zap.Logger,
	baseURL string,
	shorterService ShortedService,
	shorterRepository repositories.ShortURLRepository,
	fansShortService *fans.FansShortService,
) *ShortedHandler {
	return &ShortedHandler{
		ShorterService:    shorterService,
		ShorterRepository: shorterRepository,
		baseURL:           baseURL,
		log:               log,
		fansShortService:  fansShortService,
	}
}

// Create создание короткой ссылки
//
//	@summary Создание короткой ссылки
//	@accept  plain
//	@produce plain
//	@success 201 {string} http://localhost:8080/aAUdjf
//	@failure 409 {string} message
//	@failure 500 {string} message
//	@router  / [post]
func (sh *ShortedHandler) Create(wr http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))

	if err != nil {
		sh.log.Error("error", zap.Error(err))
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	if mediaType != "text/plain" && mediaType != "application/x-gzip" {
		http.Error(wr, "bad request", http.StatusBadRequest)
		return
	}

	var body []byte

	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		if body, err = decompress(r.Body); err != nil {
			sh.log.Error("error", zap.Error(err))
			http.Error(wr, err.Error(), http.StatusInternalServerError)

			return
		}
	} else {
		body, err = io.ReadAll(r.Body)
		if err != nil {
			sh.log.Error("error", zap.Error(err))
			http.Error(wr, err.Error(), http.StatusInternalServerError)

			return
		}
	}
	defer r.Body.Close()

	_, err = url.ParseRequestURI(string(body))

	if err != nil {
		sh.log.Error("error", zap.Error(err))
		http.Error(wr, "Invalid url", http.StatusBadRequest)
		return
	}

	userID, _ := middlewares.GetUserIDCtx(r.Context())

	shortURL, err := sh.ShorterService.Create(r.Context(), userID, string(body))

	var shortURLCreateConflictError *sdb.ShortURLCreateConflictError

	if errors.As(err, &shortURLCreateConflictError) {
		wr.WriteHeader(http.StatusConflict)
		wr.Write([]byte(fmt.Sprintf("%s/%s", sh.baseURL, shortURLCreateConflictError.OriginID)))

		return
	}

	if err != nil {
		sh.log.Error("error", zap.Error(err))
		http.Error(wr, err.Error(), http.StatusInternalServerError)

		return
	}

	wr.WriteHeader(http.StatusCreated)
	wr.Write([]byte(fmt.Sprintf("%s/%s", sh.baseURL, shortURL.ID)))
}

// Get Редирект по короткой ссылке
//
//	@summary Редирект по короткой ссылке
//	@param   id path string true "URL ID"
//	@success 307
//	@failure 404 {string} message
//	@failure 409 {string} message Was deleted
//	@router  /{id} [get]
func (sh *ShortedHandler) Get(wr http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "id")

	shortURL, ok := sh.ShorterService.GetByID(r.Context(), shortCode)

	if !ok {
		http.Error(wr, "Not Found", http.StatusNotFound)
		return
	}

	if shortURL.IsDeleted {
		http.Error(wr, "Was deleted", http.StatusGone)
		return
	}

	http.Redirect(wr, r, shortURL.URL, http.StatusTemporaryRedirect)
}

// ShortedCreateDTO data transfer object for request
type ShortedCreateDTO struct {
	URL string `json:"url" example:"https://ya.ru"`
}

// ShortedCreateDTOPool pool dto for requests
type ShortedCreateDTOPool struct {
	pool.Pool[ShortedCreateDTO]
}

// Put return object to pool
func (p *ShortedCreateDTOPool) Put(v *ShortedCreateDTO) {
	v.URL = ""
	p.Pool.Put(v)
}

var shortedCreateDTOPool = &ShortedCreateDTOPool{}

// ShortedResponseDTO data transfer object for response
type ShortedResponseDTO struct {
	Result string `json:"result" example:"http://localhost:8080/Jndshf"`
}

// ShortedResponseDTOPool pool dto for requests
type ShortedResponseDTOPool struct {
	pool.Pool[ShortedResponseDTO]
}

// Put return object to pool
func (p *ShortedResponseDTOPool) Put(v *ShortedResponseDTO) {
	v.Result = ""
	p.Pool.Put(v)
}

var shortedResponseDTOPool = &ShortedResponseDTOPool{}

// APICreate REST API обработчик создания коротких ссылок
//
//	@summary Создание короткой ссылки
//	@tags    apiShorten
//	@accept  json
//	@produce json
//	@param   request body     ShortedCreateDTO true "Ссылка для сокращения"
//	@success 201     {object} ShortedResponseDTO
//	@failure 409     {object} ShortedResponseDTO Ранее созданная короткая ссылка
//	@failure 400     {string} string             message
//	@failure 500     {string} string             message
//	@router  /api/shorten/ [post]
func (sh *ShortedHandler) APICreate(wr http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))

	if err != nil {
		sh.log.Error("error", zap.Error(err))
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	if mediaType != contentTypeJSON {
		http.Error(wr, "bad request", http.StatusBadRequest)
		return
	}

	acceptHeader := r.Header.Get("Accept")

	if acceptHeader != "" {
		crossAccepting, errAccept := accept.Negotiate(acceptHeader, contentTypeJSON)

		if errAccept != nil {
			http.Error(wr, "bad headers", http.StatusBadRequest)
			return
		}

		if crossAccepting != contentTypeJSON {
			http.Error(wr, "bad accepting content", http.StatusNotAcceptable)
			return
		}
	}

	var body []byte

	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		if body, err = decompress(r.Body); err != nil {
			sh.log.Error("error", zap.Error(err))
			http.Error(wr, err.Error(), http.StatusInternalServerError)

			return
		}
	} else {
		body, err = io.ReadAll(r.Body)
		if err != nil {
			sh.log.Error("error", zap.Error(err))
			http.Error(wr, err.Error(), http.StatusInternalServerError)

			return
		}
	}
	defer r.Body.Close()

	shortedCreateDTO := shortedCreateDTOPool.Get()
	defer shortedCreateDTOPool.Put(shortedCreateDTO)

	err = json.Unmarshal(body, &shortedCreateDTO)
	if err != nil {
		http.Error(wr, "Error parse body", http.StatusInternalServerError)
		return
	}

	_, err = url.ParseRequestURI(shortedCreateDTO.URL)
	if err != nil {
		http.Error(wr, "Invalid url", http.StatusBadRequest)
		return
	}

	userID, _ := middlewares.GetUserIDCtx(r.Context())
	shortURL, err := sh.ShorterService.Create(r.Context(), userID, shortedCreateDTO.URL)

	// TODO: Отрефакторить и убрать дублирование кода
	var shortURLCreateConflictError *sdb.ShortURLCreateConflictError

	if errors.As(err, &shortURLCreateConflictError) {
		resultURL := fmt.Sprintf("%s/%s", sh.baseURL, shortURLCreateConflictError.OriginID)

		responseCreateDTO := ShortedResponseDTO{Result: resultURL}

		responseBody, errJSONMarshal := json.Marshal(responseCreateDTO)

		if errJSONMarshal != nil {
			http.Error(wr, "error create response", http.StatusInternalServerError)
			return
		}

		wr.Header().Add("Content-Type", "application/json")
		wr.WriteHeader(http.StatusConflict)

		wr.Write(responseBody)

		return
	}

	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)

		return
	}

	resultURL := fmt.Sprintf("%s/%s", sh.baseURL, shortURL.ID)

	responseCreateDTO := shortedResponseDTOPool.Get()
	responseCreateDTO.Result = resultURL
	defer shortedResponseDTOPool.Put(responseCreateDTO)

	responseBody, err := json.Marshal(responseCreateDTO)

	if err != nil {
		http.Error(wr, "error create response", http.StatusInternalServerError)
		return
	}

	wr.Header().Add("Content-Type", "application/json")
	wr.WriteHeader(http.StatusCreated)

	wr.Write(responseBody)
}

// ShortedCreateBatchDTO data transfer object for request
type ShortedCreateBatchDTO struct {
	CorrelationID string `json:"correlation_id" example:"1"`
	OriginalURL   string `json:"original_url" example:"https://ya.ru"`
}

// ShortedResponseBatchDTO data transfer object for response
type ShortedResponseBatchDTO struct {
	CorrelationID string `json:"correlation_id" example:"1"`
	ShortURL      string `json:"short_url" example:"http://localhost:8080/JfnfgyS"`
}

// APICreateBatch Создание короткой ссылки по массиву
//
//	@summary Создание короткой ссылки по массиву
//	@tags    apiShorten
//	@accept  json
//	@produce json
//	@param   request body     []ShortedCreateBatchDTO true "Ссылки для сокращения"
//	@success 201     {array}  ShortedResponseBatchDTO
//	@failure 400     {string} string message
//	@failure 500     {string} string message
//	@router  /api/shorten/batch [post]
func (sh *ShortedHandler) APICreateBatch(wr http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))

	if err != nil {
		sh.log.Error("error", zap.Error(err))
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	if mediaType != contentTypeJSON {
		http.Error(wr, "bad request", http.StatusBadRequest)
		return
	}

	acceptHeader := r.Header.Get("Accept")

	if acceptHeader != "" {
		crossAccepting, errAccept := accept.Negotiate(acceptHeader, contentTypeJSON)

		if errAccept != nil {
			http.Error(wr, "bad headers", http.StatusBadRequest)
			return
		}

		if crossAccepting != contentTypeJSON {
			http.Error(wr, "bad accepting content", http.StatusNotAcceptable)
			return
		}
	}

	var body []byte

	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		if body, err = decompress(r.Body); err != nil {
			sh.log.Error("error", zap.Error(err))
			http.Error(wr, err.Error(), http.StatusInternalServerError)

			return
		}
	} else {
		body, err = io.ReadAll(r.Body)
		if err != nil {
			sh.log.Error("error", zap.Error(err))
			http.Error(wr, err.Error(), http.StatusInternalServerError)

			return
		}
	}
	defer r.Body.Close()

	var shortedCreateBatchDTO []*ShortedCreateBatchDTO

	err = json.Unmarshal(body, &shortedCreateBatchDTO)
	if err != nil {
		http.Error(wr, "Error parse body", http.StatusInternalServerError)
		return
	}

	userID, _ := middlewares.GetUserIDCtx(r.Context())

	shoredURLs := make([]*core.ShortURL, len(shortedCreateBatchDTO))

	for i, v := range shortedCreateBatchDTO {
		_, err = url.ParseRequestURI(v.OriginalURL)
		if err != nil {
			http.Error(wr, "Invalid url", http.StatusBadRequest)
			return
		}

		shoredURLs[i] = &core.ShortURL{
			UserID: sql.NullString{
				String: userID,
				Valid:  userID != "",
			},
			URL:           v.OriginalURL,
			CorrelationID: v.CorrelationID,
		}
	}

	if err = sh.ShorterService.CreateBatch(r.Context(), &shoredURLs); err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	resultShortURLs := make([]ShortedResponseBatchDTO, len(shoredURLs))

	for i, v := range shoredURLs {
		resultURL := fmt.Sprintf("%s/%s", sh.baseURL, v.ID)
		resultShortURLs[i] = ShortedResponseBatchDTO{ShortURL: resultURL, CorrelationID: v.CorrelationID}
	}

	responseBody, err := json.Marshal(resultShortURLs)

	if err != nil {
		http.Error(wr, "error create response", http.StatusInternalServerError)
		return
	}

	wr.Header().Add("Content-Type", "application/json")
	wr.WriteHeader(http.StatusCreated)

	wr.Write(responseBody)
}

// ShortedAllUserUResponseDTO data transfer object for response
type ShortedAllUserUResponseDTO struct {
	ShortURL    string `json:"short_url" example:"http://localhost:8080/Sjfnwf"`
	OriginalURL string `json:"original_url" example:"https://ya.ru"`
}

// APIUserURLs Получить всех коротких ссылок пользователя
//
//	@summary Получить всех коротких ссылок пользователя
//	@tags    apiShorten
//	@produce json
//	@success 200 {array} ShortedAllUserUResponseDTO
//	@success 204
//	@failure 403
//	@failure 500
//	@router  /api/user/urls [get]
func (sh *ShortedHandler) APIUserURLs(wr http.ResponseWriter, r *http.Request) {
	userID, _ := middlewares.GetUserIDCtx(r.Context())

	content, err := sh.ShorterService.AllByUser(r.Context(), userID)

	if err != nil {
		http.Error(wr, "error create response", http.StatusInternalServerError)
		return
	}

	if len(content) == 0 {
		wr.WriteHeader(http.StatusNoContent)
		return
	}

	responseDTO := make([]ShortedAllUserUResponseDTO, len(content))

	for i, shortURL := range content {
		responseDTO[i] = ShortedAllUserUResponseDTO{ShortURL: fmt.Sprintf("%s/%s", sh.baseURL, shortURL.ID), OriginalURL: shortURL.URL}
	}

	newContent, err := json.Marshal(responseDTO)

	if err != nil {
		http.Error(wr, "error create response", http.StatusInternalServerError)
		return
	}

	wr.Header().Add("Content-Type", "application/json")
	wr.Write(newContent)
}

// APIUserDeleteURLs Удаление ссылок пользователем
//
//	@summary Удаление ссылок пользователем
//	@tags    apiShorten
//	@accept  json
//	@param   request body []string true "Массив идентификаторов коротких ссылок"
//	@success 202
//	@failure 400
//	@failure 403
//	@failure 500
//	@router  /api/user/urls [delete]
func (sh *ShortedHandler) APIUserDeleteURLs(wr http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))

	if err != nil {
		sh.log.Error("error", zap.Error(err))
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	if mediaType != contentTypeJSON {
		http.Error(wr, "bad request", http.StatusBadRequest)
		return
	}

	acceptHeader := r.Header.Get("Accept")

	if acceptHeader != "" {
		crossAccepting, errAccept := accept.Negotiate(acceptHeader, contentTypeJSON)

		if errAccept != nil {
			http.Error(wr, "bad headers", http.StatusBadRequest)
			return
		}

		if crossAccepting != contentTypeJSON {
			http.Error(wr, "bad accepting content", http.StatusNotAcceptable)
			return
		}
	}

	var body []byte

	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		if body, err = decompress(r.Body); err != nil {
			sh.log.Error("error", zap.Error(err))
			http.Error(wr, err.Error(), http.StatusInternalServerError)

			return
		}
	} else {
		body, err = io.ReadAll(r.Body)
		if err != nil {
			sh.log.Error("error", zap.Error(err))
			http.Error(wr, err.Error(), http.StatusInternalServerError)

			return
		}
	}

	defer r.Body.Close()

	var urlIDs []string

	if err := json.Unmarshal(body, &urlIDs); err != nil {
		http.Error(wr, "Error parse body", http.StatusInternalServerError)
		return
	}

	userID, _ := middlewares.GetUserIDCtx(r.Context())

	sh.log.Info("was delete", zap.String("userID", userID), zap.Strings("urlIDs", urlIDs))

	sh.fansShortService.Add(userID, urlIDs)

	wr.WriteHeader(http.StatusAccepted)
}

func decompress(dateRead io.Reader) ([]byte, error) {
	gr, err := gzip.NewReader(dateRead)

	if err != nil {
		return nil, err
	}

	defer gr.Close()

	var b bytes.Buffer

	if _, err := b.ReadFrom(gr); err != nil {
		return nil, fmt.Errorf("failed decopress data :%w", err)
	}

	return b.Bytes(), nil
}

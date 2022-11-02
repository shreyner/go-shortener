package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	service2 "github.com/shreyner/go-shortener/internal/service"
	"github.com/shreyner/go-shortener/internal/storage"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/shreyner/go-shortener/internal/core"
)

// TODO: Проверять сообщения при плохих ответах

var (
	ContentType = "text/plain; charset=utf-8"
)

type MyMockService struct {
	mock.Mock
}

func (m *MyMockService) Create(_ context.Context, id, url string) (*core.ShortURL, error) {
	args := m.Called(id, url)
	return args.Get(0).(*core.ShortURL), args.Error(1)
}

func (m *MyMockService) GetByID(_ context.Context, key string) (*core.ShortURL, bool) {
	args := m.Called(key)
	return args.Get(0).(*core.ShortURL), args.Bool(1)
}

func (m *MyMockService) AllByUser(_ context.Context, id string) ([]*core.ShortURL, error) {
	args := m.Called(id)

	shortURLs, ok := args.Get(0).([]*core.ShortURL)
	if !ok {
		log.Print("Error in type")
	}

	return shortURLs, args.Error(1)
}

func (m *MyMockService) CreateBatch(ctx context.Context, shortURLs *[]*core.ShortURL) error {
	args := m.Called(ctx, shortURLs)
	return args.Error(0)
}

func testRequest(t *testing.T, ts *httptest.Server, method, path, contentType, accept, body string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBufferString(body))

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	if accept != "" {
		req.Header.Set("Accept", accept)
	}

	require.NoError(t, err)

	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("Some Error", err)
	}

	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestShortedHandler_ShortedCreate(t *testing.T) {
	t.Run("should success create", func(t *testing.T) {
		mockService := new(MyMockService)

		r := NewRouter(zap.NewNop(), "http://localhost:8080", mockService, nil, nil, nil)
		ts := httptest.NewServer(r)

		mockService.On("Create", mock.Anything, "https://ya.ru/").Return(
			&core.ShortURL{URL: "https://ya.ru/", ID: "ya"},
			nil,
		)

		resp, respBody := testRequest(t, ts, http.MethodPost, "/", ContentType, "", "https://ya.ru/")
		defer resp.Body.Close()

		mockService.AssertExpectations(t)
		mockService.AssertCalled(t, "Create", mock.AnythingOfType("string"), "https://ya.ru/")
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, "http://localhost:8080/ya", respBody)
	})

	t.Run("should error for incorrect url", func(t *testing.T) {
		mockService := new(MyMockService)

		r := NewRouter(zap.NewNop(), "http://localhost:8080", mockService, nil, nil, nil)
		ts := httptest.NewServer(r)

		resp, _ := testRequest(t, ts, http.MethodPost, "/", ContentType, "", "yyy")
		defer resp.Body.Close()

		mockService.AssertNotCalled(t, "Create")
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should error for incorrect Content-Type", func(t *testing.T) {
		mockService := new(MyMockService)

		r := NewRouter(zap.NewNop(), "http://localhost:8080", mockService, nil, nil, nil)
		ts := httptest.NewServer(r)

		resp, _ := testRequest(t, ts, http.MethodPost, "/", "text/html; charset=utf8", "", "https://ya.ru/")
		defer resp.Body.Close()

		mockService.AssertNotCalled(t, "Create")
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should error method not allowed for POST /some", func(t *testing.T) {
		mockService := new(MyMockService)

		r := NewRouter(zap.NewNop(), "http://localhost:8080", mockService, nil, nil, nil)
		ts := httptest.NewServer(r)

		resp, _ := testRequest(t, ts, http.MethodPost, "/some", ContentType, "", "https://ya.ru/")
		defer resp.Body.Close()

		mockService.AssertNotCalled(t, "Create")
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	// TODO: Fixed
	//t.Run("should error for incorrect without content-type", func(t *testing.T) {
	//	mockService := new(MyMockService)
	//
	//	r := NewRouter(zap.NewNop(),"http://localhost:8080", mockService)
	//	ts := httptest.NewServer(r)
	//
	//	resp, _ := testRequest(t, ts, http.MethodPost, "/", "", "https://ya.ru/")
	//  defer resp.Body.Close()
	//
	//	mockService.AssertNotCalled(t, "Create")
	//	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	//})
}

func TestShortedHandler_ShortedGet(t *testing.T) {
	t.Run("should success redirect", func(t *testing.T) {
		mockService := new(MyMockService)

		r := NewRouter(zap.NewNop(), "http://localhost:8080", mockService, nil, nil, nil)
		ts := httptest.NewServer(r)

		mockService.On("GetByID", "asdd").Return(&core.ShortURL{ID: "asdd", URL: "https://ya.ru"}, true)

		resp, _ := testRequest(t, ts, http.MethodGet, "/asdd", "", "", "")
		defer resp.Body.Close()

		mockService.AssertExpectations(t)
		mockService.AssertCalled(t, "GetByID", "asdd")
		require.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
		assert.Equal(t, "https://ya.ru", resp.Header.Get("Location"))
	})

	t.Run("should error for not found by id", func(t *testing.T) {
		mockService := new(MyMockService)

		r := NewRouter(zap.NewNop(), "http://localhost:8080", mockService, nil, nil, nil)
		ts := httptest.NewServer(r)

		mockService.On("GetByID", "not").Return(&core.ShortURL{}, false)

		resp, _ := testRequest(t, ts, http.MethodGet, "/not", "", "", "")
		defer resp.Body.Close()

		mockService.AssertExpectations(t)
		mockService.AssertCalled(t, "GetByID", "not")
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestShortedHandler_ApiCreate(t *testing.T) {
	t.Run("should success create", func(t *testing.T) {
		contentType := "application/json"
		acceptType := "application/json"
		mockService := new(MyMockService)

		r := NewRouter(zap.NewNop(), "http://localhost:8080", mockService, nil, nil, nil)
		ts := httptest.NewServer(r)

		mockService.On("Create", mock.Anything, "https://ya.ru/").Return(&core.ShortURL{URL: "https://ya.ru/", ID: "ya"}, nil)

		resp, respBody := testRequest(t, ts, http.MethodPost, "/api/shorten", contentType, acceptType, `{"url":"https://ya.ru/"}`)
		defer resp.Body.Close()

		mockService.AssertExpectations(t)
		mockService.AssertCalled(t, "Create", mock.AnythingOfType("string"), "https://ya.ru/")
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, "{\"result\":\"http://localhost:8080/ya\"}", respBody)
	})

	t.Run("should return currect ContentType", func(t *testing.T) {
		contentType := "application/json"
		acceptType := "application/json"
		mockService := new(MyMockService)

		r := NewRouter(zap.NewNop(), "http://localhost:8080", mockService, nil, nil, nil)
		ts := httptest.NewServer(r)

		mockService.On("Create", mock.Anything, "https://ya.ru/").Return(&core.ShortURL{URL: "https://ya.ru/", ID: "ya"}, nil)

		resp, respBody := testRequest(t, ts, http.MethodPost, "/api/shorten", contentType, acceptType, "{\"url\":\"https://ya.ru/\"}")
		defer resp.Body.Close()

		mockService.AssertExpectations(t)
		mockService.AssertCalled(t, "Create", mock.AnythingOfType("string"), "https://ya.ru/")
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, "{\"result\":\"http://localhost:8080/ya\"}", respBody)
	})

	t.Run("should error for incorrect url", func(t *testing.T) {
		contentType := "application/json"
		acceptType := "application/json"
		mockService := new(MyMockService)

		r := NewRouter(zap.NewNop(), "http://localhost:8080", mockService, nil, nil, nil)
		ts := httptest.NewServer(r)

		mockService.On("Create", mock.Anything, "https://ya.ru/").Return(&core.ShortURL{URL: "https://ya.ru/", ID: "ya"}, nil)

		resp, _ := testRequest(t, ts, http.MethodPost, "/api/shorten", contentType, acceptType, "{\"url\":\"ya\"}")
		defer resp.Body.Close()

		mockService.AssertNotCalled(t, "Create")
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should error for incorrect Content-Type", func(t *testing.T) {
		acceptType := "application/json"
		mockService := new(MyMockService)

		r := NewRouter(zap.NewNop(), "http://localhost:8080", mockService, nil, nil, nil)
		ts := httptest.NewServer(r)

		resp, _ := testRequest(t, ts, http.MethodPost, "/api/shorten", "text/plain", acceptType, "{\"url\":\"https://ya.ru/\"}")
		defer resp.Body.Close()

		mockService.AssertNotCalled(t, "Create")
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should error for incorrect Accept", func(t *testing.T) {
		contentType := "application/json"
		mockService := new(MyMockService)

		r := NewRouter(zap.NewNop(), "http://localhost:8080", mockService, nil, nil, nil)
		ts := httptest.NewServer(r)

		resp, _ := testRequest(t, ts, http.MethodPost, "/api/shorten", contentType, "application/xml", "{\"url\":\"https://ya.ru/\"}")
		defer resp.Body.Close()

		mockService.AssertNotCalled(t, "Create")
		assert.Equal(t, http.StatusNotAcceptable, resp.StatusCode)
	})
}

func BenchmarkShortedHandler_APICreate(b *testing.B) {
	var indexRequest int64 = 0
	memoRepository, _ := storage.NewStorage(
		zap.NewNop(),
		"",
		"",
	)
	defer memoRepository.Close()

	service := service2.NewService(
		memoRepository.ShortURL,
	)

	shortedHandler := NewShortedHandler(
		zap.NewNop(),
		"http://localhost:8080",
		service.ShorterService,
		memoRepository.ShortURL,
		nil,
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		index := atomic.AddInt64(&indexRequest, 1)
		body := fmt.Sprintf(`{"url": "https://ya.ru/%v"}`, index)
		req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/shorten/", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		wr := httptest.NewRecorder()

		b.StartTimer()
		shortedHandler.APICreate(wr, req)
		b.StopTimer()

		if wr.Code != http.StatusCreated {
			b.Fatal("unexpected response status code", wr.Code)
		}
	}
}

func ExampleShortedHandler_Create() {
	req, err := http.NewRequest("POST", "http://localhost:8080/", strings.NewReader("https://yandex.ru"))

	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "text/plain")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return
	}

	fmt.Println(body) // http://localhost:8080/Sfnvdf
}

func ExampleShortedHandler_Get() {
	req, err := http.NewRequest("GET", "http://localhost:8080/Sfnvdf", nil)

	if err != nil {
		return
	}

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Do(req)

	if err != nil {
		return
	}

	defer res.Body.Close()

	URL := res.Header.Get("Location")

	fmt.Println(URL) // https://ya.ru
}

func ExampleShortedHandler_APICreate() {
	requestBody := `
		{
			"url": "https://ya.ru"
		}
	`
	userCookie := &http.Cookie{
		Name:  "auth",
		Value: "authCookieFormFirstRequest",
		Path:  "/",
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/api/shorten/", strings.NewReader(requestBody))

	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(userCookie)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return
	}

	defer res.Body.Close()

	bodyRaw, err := io.ReadAll(res.Body)

	if err != nil {
		return
	}

	var result ShortedResponseDTO
	if err = json.Unmarshal(bodyRaw, &result); err != nil {
		return
	}

	fmt.Println(result.Result) // http://localhost:8080/Sfnvdf
}

func ExampleShortedHandler_APICreateBatch() {
	requestBody := `
		[
			{
				"correlation_id": "1",
				"original_url": "https://ya.ru"
			},
			{
				"correlation_id": "2",
				"original_url": "https://vk.com"
			},
		]
	`
	userCookie := &http.Cookie{
		Name:  "auth",
		Value: "authCookieFormFirstRequest",
		Path:  "/",
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/api/shorten/batch", strings.NewReader(requestBody))

	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(userCookie)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return
	}

	defer res.Body.Close()

	bodyRaw, err := io.ReadAll(res.Body)

	if err != nil {
		return
	}

	var result []*ShortedResponseBatchDTO
	err = json.Unmarshal(bodyRaw, &result)

	if err != nil {
		return
	}

	for _, responseDTO := range result {
		fmt.Printf("correlation: %v, short: %v", responseDTO.CorrelationID, responseDTO.ShortURL)
	}

	// correlation: 1, short: http://localhost:8080/aCewfns
	// correlation: 2, short: http://localhost:8080/aMJFSNjs
}

func ExampleShortedHandler_APIUserURLs() {
	userCookie := &http.Cookie{
		Name:  "auth",
		Value: "authCookieFormFirstRequest",
		Path:  "/",
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/api/user/urls", nil)

	if err != nil {
		return
	}

	req.AddCookie(userCookie)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return
	}

	defer res.Body.Close()

	if res.StatusCode == 204 {
		fmt.Println("User hasn't shorten URL's")
		return
	}

	bodyRaw, err := io.ReadAll(res.Body)

	if err != nil {
		return
	}

	var result []*ShortedAllUserUResponseDTO
	err = json.Unmarshal(bodyRaw, &result)

	if err != nil {
		return
	}

	for _, responseDTO := range result {
		fmt.Printf("original: %v, short: %v", responseDTO.OriginalURL, responseDTO.ShortURL)
	}

	// original: https://ya.ru, short: http://localhost:8080/aCewfns
	// original: https://vk.com, short: http://localhost:8080/aMJFSNjs
}

func ExampleShortedHandler_APIUserDeleteURLs() {
	requestBody := `["aCewfns", "aMJFSNjs"]`
	userCookie := &http.Cookie{
		Name:  "auth",
		Value: "authCookieFormFirstRequest",
		Path:  "/",
	}

	req, err := http.NewRequest("DELETE", "http://localhost:8080/api/user/urls", strings.NewReader(requestBody))

	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(userCookie)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return
	}

	defer res.Body.Close()

	if res.StatusCode != 202 {
		fmt.Println("Error response")
		return
	}
}

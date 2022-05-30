package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shreyner/go-shortener/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TODO: Проверять сообщения при плохих ответах

var (
	ContentType = "application/x-www-form-urlencoded"
)

type MyMockService struct {
	mock.Mock
}

func (m *MyMockService) Create(url string) (core.ShortUrl, error) {
	args := m.Called(url)
	return args.Get(0).(core.ShortUrl), args.Error(1)
}

func (m *MyMockService) GetById(key string) (core.ShortUrl, bool) {
	args := m.Called(key)
	return args.Get(0).(core.ShortUrl), args.Bool(1)
}

func testRequest(t *testing.T, ts *httptest.Server, method, path, contentType, body string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBufferString(body))

	if method == http.MethodPost && len(contentType) > 0 {
		req.Header.Set("Content-Type", contentType)
	}

	require.NoError(t, err)

	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestShortedHandler_ShortedCreate(t *testing.T) {
	t.Run("should success create", func(t *testing.T) {
		mockService := new(MyMockService)

		r := NewRouter(mockService)
		ts := httptest.NewServer(r)

		mockService.On("Create", "https://ya.ru/").Return(core.ShortUrl{Url: "https://ya.ru/", Id: "ya"}, nil)

		resp, body := testRequest(t, ts, http.MethodPost, "/", ContentType, "https://ya.ru/")

		mockService.AssertExpectations(t)
		mockService.AssertCalled(t, "Create", "https://ya.ru/")
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, "ya", body)
	})

	t.Run("should error for incorrect url", func(t *testing.T) {
		mockService := new(MyMockService)

		r := NewRouter(mockService)
		ts := httptest.NewServer(r)

		resp, _ := testRequest(t, ts, http.MethodPost, "/", ContentType, "yyy")

		mockService.AssertNotCalled(t, "Create")
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should error for incorrect Content-Type", func(t *testing.T) {
		mockService := new(MyMockService)

		r := NewRouter(mockService)
		ts := httptest.NewServer(r)

		resp, _ := testRequest(t, ts, http.MethodPost, "/", "text/html; charset=utf8", "https://ya.ru/")

		mockService.AssertNotCalled(t, "Create")
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should error method not allowed for POST /some", func(t *testing.T) {
		mockService := new(MyMockService)

		r := NewRouter(mockService)
		ts := httptest.NewServer(r)

		resp, _ := testRequest(t, ts, http.MethodPost, "/some", ContentType, "https://ya.ru/")

		mockService.AssertNotCalled(t, "Create")
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	// TODO: Fixed
	//t.Run("should error for incorrect without content-type", func(t *testing.T) {
	//	mockService := new(MyMockService)
	//
	//	r := NewRouter(mockService)
	//	ts := httptest.NewServer(r)
	//
	//	resp, _ := testRequest(t, ts, http.MethodPost, "/", "", "https://ya.ru/")
	//
	//	mockService.AssertNotCalled(t, "Create")
	//	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	//})
}

func TestShortedHandler_ShortedGet(t *testing.T) {
	t.Run("should success redirect", func(t *testing.T) {
		mockService := new(MyMockService)

		r := NewRouter(mockService)
		ts := httptest.NewServer(r)

		mockService.On("GetById", "asdd").Return(core.ShortUrl{Id: "asdd", Url: "https://ya.ru"}, true)

		resp, _ := testRequest(t, ts, http.MethodGet, "/asdd", "", "")

		mockService.AssertExpectations(t)
		mockService.AssertCalled(t, "GetById", "asdd")
		require.Equal(t, http.StatusPermanentRedirect, resp.StatusCode)
		assert.Equal(t, "https://ya.ru", resp.Header.Get("Location"))
	})

	t.Run("should error for not found by id", func(t *testing.T) {
		mockService := new(MyMockService)

		r := NewRouter(mockService)
		ts := httptest.NewServer(r)

		mockService.On("GetById", "not").Return(core.ShortUrl{}, false)

		resp, _ := testRequest(t, ts, http.MethodGet, "/not", "", "")

		mockService.AssertExpectations(t)
		mockService.AssertCalled(t, "GetById", "not")
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

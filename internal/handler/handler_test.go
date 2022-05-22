package handler_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shreyner/go-shortener/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexHandler(t *testing.T) {

	t.Run("shuold create link", func(t *testing.T) {
		storeMock := map[string]string{}
		request := httptest.NewRequest("POST", "/", bytes.NewBufferString("https://yandex.ru"))
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		h := handler.IndexHandler(storeMock)

		h.ServeHTTP(w, request)
		result := w.Result()

		require.Equal(t, http.StatusCreated, result.StatusCode)

		bodyResult, err := ioutil.ReadAll(result.Body)
		require.NoError(t, err)
		err = result.Body.Close()
		require.NoError(t, err)

		strBody := string(bodyResult[:])

		targetUrl, ok := storeMock[strBody]

		require.Equal(t, true, ok)
		assert.Equal(t, "https://yandex.ru", targetUrl)
	})

	t.Run("shuold error for incorrect url", func(t *testing.T) {
		storeMock := map[string]string{}
		request := httptest.NewRequest("POST", "/", bytes.NewBufferString("yandex"))
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		h := handler.IndexHandler(storeMock)

		h.ServeHTTP(w, request)
		result := w.Result()

		require.Equal(t, http.StatusBadRequest, result.StatusCode)
	})

	t.Run("shuold error for incorrect Content-Type", func(t *testing.T) {
		storeMock := map[string]string{}
		request := httptest.NewRequest("POST", "/", bytes.NewBufferString("yandex"))
		request.Header.Add("Content-Type", "text/plain")

		w := httptest.NewRecorder()
		h := handler.IndexHandler(storeMock)

		h.ServeHTTP(w, request)
		result := w.Result()

		require.Equal(t, http.StatusBadRequest, result.StatusCode)
	})

	t.Run("shuold error for incorrect method for index path", func(t *testing.T) {
		storeMock := map[string]string{}
		request := httptest.NewRequest("GET", "/", nil)

		w := httptest.NewRecorder()
		h := handler.IndexHandler(storeMock)

		h.ServeHTTP(w, request)
		result := w.Result()

		require.Equal(t, http.StatusMethodNotAllowed, result.StatusCode)
	})

	t.Run("shuold error not found for GET", func(t *testing.T) {
		storeMock := map[string]string{}
		request := httptest.NewRequest("GET", "/some/invalid", nil)

		w := httptest.NewRecorder()
		h := handler.IndexHandler(storeMock)

		h.ServeHTTP(w, request)
		result := w.Result()

		require.Equal(t, http.StatusNotFound, result.StatusCode)
	})

	t.Run("shuold error not found for POST", func(t *testing.T) {
		storeMock := map[string]string{}
		request := httptest.NewRequest("POST", "/some/invalid", nil)

		w := httptest.NewRecorder()
		h := handler.IndexHandler(storeMock)

		h.ServeHTTP(w, request)
		result := w.Result()

		require.Equal(t, http.StatusNotFound, result.StatusCode)
	})

	t.Run("shuold error method not allowed for POST /some", func(t *testing.T) {
		storeMock := map[string]string{}
		request := httptest.NewRequest("POST", "/some", nil)

		w := httptest.NewRecorder()
		h := handler.IndexHandler(storeMock)

		h.ServeHTTP(w, request)
		result := w.Result()

		require.Equal(t, http.StatusMethodNotAllowed, result.StatusCode)
	})

	t.Run("shuold success redirect", func(t *testing.T) {
		storeMock := map[string]string{
			"yand": "https://yandex.ru",
		}
		request := httptest.NewRequest("GET", "/yand", nil)

		w := httptest.NewRecorder()
		h := handler.IndexHandler(storeMock)

		h.ServeHTTP(w, request)
		result := w.Result()

		require.Equal(t, http.StatusPermanentRedirect, result.StatusCode)
		assert.Equal(t, "https://yandex.ru", result.Header.Get("Location"))
	})

	t.Run("shuold error for not found by id", func(t *testing.T) {
		storeMock := map[string]string{}
		request := httptest.NewRequest("GET", "/TGss", nil)

		w := httptest.NewRecorder()
		h := handler.IndexHandler(storeMock)

		h.ServeHTTP(w, request)
		result := w.Result()

		require.Equal(t, http.StatusNotFound, result.StatusCode)
	})

}

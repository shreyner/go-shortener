package handlers

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO: ВА чем разница между io и ioutil пкетов

//type MyMockService {
//	mock
//}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestShortedHandler_ShortedGet(t *testing.T) {
	t.Run("should success redirect", func(t *testing.T) {

		r := NewRouter()

	})

	//type fields struct {
	//	ShorterService ShortedService
	//}
	//type args struct {
	//	wr http.ResponseWriter
	//	r  *http.Request
	//}
	//tests := []struct {
	//	name   string
	//	fields fields
	//	args   args
	//}{
	//	// TODO: Add test cases.
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		sh := &ShortedHandler{
	//			ShorterService: tt.fields.ShorterService,
	//		}
	//		sh.ShortedGet(tt.args.wr, tt.args.r)
	//	})
	//}
}

//func TestShortedHandler_ShortedCreate(t *testing.T) {
//	type fields struct {
//		ShorterService ShortedService
//	}
//	type args struct {
//		wr http.ResponseWriter
//		r  *http.Request
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			sh := &ShortedHandler{
//				ShorterService: tt.fields.ShorterService,
//			}
//			sh.ShortedCreate(tt.args.wr, tt.args.r)
//		})
//	}
//}

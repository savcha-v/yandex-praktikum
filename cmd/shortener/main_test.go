package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShort(t *testing.T) {
	type want struct {
		body   string
		status int
	}
	tests := []struct {
		name    string
		request string
		method  string
		body    io.Reader
		want    want
	}{
		{
			name:    "test endpoint GET",
			method:  http.MethodGet,
			request: "/?id=1",
			want: want{
				status: http.StatusBadRequest,
				body:   "'id' not found\n",
			},
		},
		{
			name:    "test endpoint DELETE",
			method:  http.MethodDelete,
			request: "/?id=1",
			want: want{
				status: http.StatusMethodNotAllowed,
				body:   "POST or GET requests are allowed!\n",
			},
		},
		{
			name:    "test endpoint PUT",
			method:  http.MethodPut,
			request: "/?id=1",
			want: want{
				status: http.StatusMethodNotAllowed,
				body:   "POST or GET requests are allowed!\n",
			},
		},
		{
			name:    "test endpoint PATCH",
			method:  http.MethodPatch,
			request: "/?id=1",
			want: want{
				status: http.StatusMethodNotAllowed,
				body:   "POST or GET requests are allowed!\n",
			},
		},
		{
			name:    "test endpoint CONNECT",
			method:  http.MethodConnect,
			request: "/?id=1",
			want: want{
				status: http.StatusMethodNotAllowed,
				body:   "POST or GET requests are allowed!\n",
			},
		},
		{
			name:    "test endpoint HEAD",
			method:  http.MethodHead,
			request: "/?id=1",
			want: want{
				status: http.StatusMethodNotAllowed,
				body:   "POST or GET requests are allowed!\n",
			},
		},
		{
			name:    "test endpoint OPTIONS",
			method:  http.MethodOptions,
			request: "/?id=1",
			want: want{
				status: http.StatusMethodNotAllowed,
				body:   "POST or GET requests are allowed!\n",
			},
		},
		{
			name:    "test endpoint TRACE",
			method:  http.MethodTrace,
			request: "/?id=1",
			want: want{
				status: http.StatusMethodNotAllowed,
				body:   "POST or GET requests are allowed!\n",
			},
		},
		{
			name:    "test endpoint POST",
			method:  http.MethodPost,
			request: "/",
			body:    strings.NewReader("https://quickref.me/golang"),
			want: want{
				body:   "http://example.com/?id=1",
				status: http.StatusCreated,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, tt.body)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(Short)
			h.ServeHTTP(w, request)
			result := w.Result()

			require.Equal(t, tt.want.status, result.StatusCode)

			bodyResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.body, string(bodyResult))

		})
	}
}

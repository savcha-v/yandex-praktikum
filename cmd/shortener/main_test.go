package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type Want struct {
	body   string
	status int
}
type TestStruct struct {
	name    string
	request string
	method  string
	body    io.Reader
	want    Want
}

var TestsForAll = []TestStruct{
	{
		name:    "test endpoint DELETE",
		method:  http.MethodDelete,
		request: "/?id=1",
		want: Want{
			status: http.StatusMethodNotAllowed,
			body:   "POST or GET requests are allowed!\n",
		},
	},
	{
		name:    "test endpoint PUT",
		method:  http.MethodPut,
		request: "/?id=1",
		want: Want{
			status: http.StatusMethodNotAllowed,
			body:   "POST or GET requests are allowed!\n",
		},
	},
	{
		name:    "test endpoint PATCH",
		method:  http.MethodPatch,
		request: "/?id=1",
		want: Want{
			status: http.StatusMethodNotAllowed,
			body:   "POST or GET requests are allowed!\n",
		},
	},
	{
		name:    "test endpoint CONNECT",
		method:  http.MethodConnect,
		request: "/?id=1",
		want: Want{
			status: http.StatusMethodNotAllowed,
			body:   "POST or GET requests are allowed!\n",
		},
	},
	{
		name:    "test endpoint HEAD",
		method:  http.MethodHead,
		request: "/?id=1",
		want: Want{
			status: http.StatusMethodNotAllowed,
			body:   "POST or GET requests are allowed!\n",
		},
	},
	{
		name:    "test endpoint OPTIONS",
		method:  http.MethodOptions,
		request: "/?id=1",
		want: Want{
			status: http.StatusMethodNotAllowed,
			body:   "POST or GET requests are allowed!\n",
		},
	},
	{
		name:    "test endpoint TRACE",
		method:  http.MethodTrace,
		request: "/?id=1",
		want: Want{
			status: http.StatusMethodNotAllowed,
			body:   "POST or GET requests are allowed!\n",
		},
	},
}

var TestsForPost = []TestStruct{
	{
		name:    "test endpoint POST",
		method:  http.MethodPost,
		request: "/",
		body:    strings.NewReader("https://quickref.me/golang"),
		want: Want{
			body:   "http://example.com/?id=1",
			status: http.StatusCreated,
		},
	},
}

var TestsForGet = []TestStruct{
	{
		name:    "test endpoint GET",
		method:  http.MethodGet,
		request: "/?id=1",
		want: Want{
			status: http.StatusBadRequest,
			body:   "'id' not found\n",
		},
	},
}

func changeResult(t *testing.T, result *http.Response, want Want) {

	require.Equal(t, want.status, result.StatusCode)

	_, err := ioutil.ReadAll(result.Body)
	require.NoError(t, err)
	err = result.Body.Close()
	require.NoError(t, err)
}

func TestGetShort(t *testing.T) {
	tests := [][]TestStruct{
		TestsForAll,
		TestsForGet,
	}
	for _, maptt := range tests {
		for _, tt := range maptt {
			t.Run(tt.name, func(t *testing.T) {
				request := httptest.NewRequest(tt.method, tt.request, tt.body)
				w := httptest.NewRecorder()
				h := http.HandlerFunc(GetShort)
				h.ServeHTTP(w, request)
				result := w.Result()
				defer result.Body.Close()
				changeResult(t, result, tt.want)
			})
		}
	}

}

func TestPostShort(t *testing.T) {

	tests := [][]TestStruct{
		TestsForAll,
		TestsForPost,
	}
	for _, maptt := range tests {
		for _, tt := range maptt {
			t.Run(tt.name, func(t *testing.T) {
				request := httptest.NewRequest(tt.method, tt.request, tt.body)
				w := httptest.NewRecorder()
				h := http.HandlerFunc(PostShort)
				h.ServeHTTP(w, request)
				result := w.Result()
				defer result.Body.Close()
				changeResult(t, result, tt.want)
			})
		}
	}
}

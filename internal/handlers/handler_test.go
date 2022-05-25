package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"yandex-praktikum/internal/config"

	"github.com/caarlos0/env/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Want struct {
	response    string
	code        int
	contentType string
}
type TestStruct struct {
	name    string
	request string
	method  string
	want    Want
}

func TestPostShort(t *testing.T) {
	var cfg config.Config
	tests := []TestStruct{
		{
			name:   "test endpoint DELETE",
			method: http.MethodDelete,
			want: Want{
				code:     http.StatusBadRequest,
				response: "Shortcut url not found\n",
			},
		},
		{
			name:   "test endpoint PUT",
			method: http.MethodPut,
			want: Want{
				code:     http.StatusBadRequest,
				response: "Shortcut url not found\n",
			},
		},
		{
			name:   "test endpoint PATCH",
			method: http.MethodPatch,
			want: Want{
				code:     http.StatusBadRequest,
				response: "Shortcut url not found\n",
			},
		},
		{
			name:   "test endpoint CONNECT",
			method: http.MethodConnect,
			want: Want{
				code:     http.StatusBadRequest,
				response: "Shortcut url not found\n",
			},
		},
		{
			name:   "test endpoint HEAD",
			method: http.MethodHead,
			want: Want{
				code:     http.StatusBadRequest,
				response: "Shortcut url not found\n",
			},
		},
		{
			name:   "test endpoint OPTIONS",
			method: http.MethodOptions,
			want: Want{
				code:     http.StatusBadRequest,
				response: "Shortcut url not found\n",
			},
		},
		{
			name:   "test endpoint TRACE",
			method: http.MethodTrace,
			want: Want{
				code:     http.StatusBadRequest,
				response: "Shortcut url not found\n",
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/status", nil)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(PostShort(cfg))

			h.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			require.Equal(t, tt.want.code, res.StatusCode)

			// получаем и проверяем тело запроса
			if tt.want.response != "" {
				defer res.Body.Close()

				resBody, err := io.ReadAll(res.Body)

				assert.NoError(t, err)
				assert.Equal(t, tt.want.response, string(resBody))
			}
			// заголовок ответа
			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func TestPostGet(t *testing.T) {
	tests := []TestStruct{
		{
			name:    "test endpoint DELETE",
			method:  http.MethodDelete,
			request: "1",
			want: Want{
				code:     http.StatusBadRequest,
				response: "'id' missing\n",
			},
		},
		{
			name:    "test endpoint PUT",
			method:  http.MethodPut,
			request: "1",
			want: Want{
				code:     http.StatusBadRequest,
				response: "'id' missing\n",
			},
		},
		{
			name:    "test endpoint PATCH",
			method:  http.MethodPatch,
			request: "1",
			want: Want{
				code:     http.StatusBadRequest,
				response: "'id' missing\n",
			},
		},
		{
			name:    "test endpoint CONNECT",
			method:  http.MethodConnect,
			request: "1",
			want: Want{
				code:     http.StatusBadRequest,
				response: "'id' missing\n",
			},
		},
		{
			name:    "test endpoint HEAD",
			method:  http.MethodHead,
			request: "1",
			want: Want{
				code:     http.StatusBadRequest,
				response: "'id' missing\n",
			},
		},
		{
			name:    "test endpoint OPTIONS",
			method:  http.MethodOptions,
			request: "1",
			want: Want{
				code:     http.StatusBadRequest,
				response: "'id' missing\n",
			},
		},
		{
			name:    "test endpoint TRACE",
			method:  http.MethodTrace,
			request: "1",
			want: Want{
				code:     http.StatusBadRequest,
				response: "'id' missing\n",
			},
		},
	}

	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/status", nil)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(GetShort(cfg))

			h.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			require.Equal(t, tt.want.code, res.StatusCode)

			// получаем и проверяем тело запроса
			if tt.want.response != "" {
				defer res.Body.Close()

				resBody, err := io.ReadAll(res.Body)

				assert.NoError(t, err)
				assert.Equal(t, tt.want.response, string(resBody))
			}
			// заголовок ответа
			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func TestPostShorten(t *testing.T) {

	type bodyReqType struct {
		URL string
	}

	type bodyResType struct {
		Result string
	}

	type wantType struct {
		bodyRes     bodyResType
		code        int
		contentType string
		response    string
	}

	type testStruct struct {
		name    string
		bodyReq bodyReqType
		method  string
		want    wantType
	}

	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	bodyReq := bodyReqType{
		URL: "https://quickref.me/golang",
	}

	bodyRes := bodyResType{
		Result: "http://example.com/" + cfg.BaseURL + "0",
	}

	tests := []testStruct{
		{
			name:   "test endpoint DELETE",
			method: http.MethodDelete,
			want: wantType{
				code:     http.StatusBadRequest,
				response: "'id' missing\n",
			},
		},
		{
			name:   "test endpoint PUT",
			method: http.MethodPut,
			want: wantType{
				code:     http.StatusBadRequest,
				response: "'id' missing\n",
			},
		},
		{
			name:   "test endpoint PATCH",
			method: http.MethodPatch,

			want: wantType{
				code:     http.StatusBadRequest,
				response: "'id' missing\n",
			},
		},
		{
			name:   "test endpoint CONNECT",
			method: http.MethodConnect,
			want: wantType{
				code:     http.StatusBadRequest,
				response: "'id' missing\n",
			},
		},
		{
			name:   "test endpoint HEAD",
			method: http.MethodHead,
			want: wantType{
				code:     http.StatusBadRequest,
				response: "'id' missing\n",
			},
		},
		{
			name:   "test endpoint OPTIONS",
			method: http.MethodOptions,
			want: wantType{
				code:     http.StatusBadRequest,
				response: "'id' missing\n",
			},
		},
		{
			name:   "test endpoint TRACE",
			method: http.MethodTrace,
			want: wantType{
				code:     http.StatusBadRequest,
				response: "'id' missing\n",
			},
		},
		{
			name:    "test endpoint POST",
			method:  http.MethodPost,
			bodyReq: bodyReq,
			want: wantType{
				code:    http.StatusCreated,
				bodyRes: bodyRes,
			},
		},
	}

	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			h := http.HandlerFunc(PostShorten(cfg))

			if tt.bodyReq != (bodyReqType{}) {
				requestBody, err := json.Marshal(&bodyReq)
				assert.NoError(t, err)
				request := httptest.NewRequest(tt.method, "/status", bytes.NewReader(requestBody))
				h.ServeHTTP(w, request)
			} else {
				request := httptest.NewRequest(tt.method, "/status", nil)
				h.ServeHTTP(w, request)
			}

			res := w.Result()
			defer res.Body.Close()

			// проверяем код ответа
			require.Equal(t, tt.want.code, res.StatusCode)

			// заголовок ответа
			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			}

			resBody, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

			if tt.want.bodyRes != (bodyResType{}) {
				bodyRes := bodyResType{}
				json.NewDecoder(bytes.NewReader(resBody)).Decode(&bodyRes)
				assert.Equal(t, tt.want.bodyRes, bodyRes)

				assert.NoError(t, err)

			}
		})
	}
}

package handlers

import (
	"github.com/IgorPestretsov/yandex_shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetFullLinkByID(t *testing.T) {
	type want struct {
		statusCode int
		location   string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "test1",
			request: "/wrong",
			want: want{
				statusCode: http.StatusBadRequest,
				location:   "",
			},
		},
		{
			name:    "test2",
			request: "/ggl",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "https://google.com",
			},
		},
		{
			name:    "test3",
			request: "/yndxprct",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "https://practicum.yandex.ru",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.New("")
			s.SaveLinksPair("https://google.com", "ggl")
			s.SaveLinksPair("https://practicum.yandex.ru", "yndxprct")

			r := chi.NewRouter()

			r.Get("/{id}", func(rw http.ResponseWriter, r *http.Request) {
				GetFullLinkByID(rw, r, s)
			})

			req := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, result.Header.Get("Location"), tt.want.location)
		})
	}
}

func TestGetShortLink(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name string
		body io.Reader
		want want
	}{
		{
			name: "test 1",
			body: strings.NewReader("yandex.ru"),
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name: "test 2",
			body: strings.NewReader(""),
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.New("")

			r := chi.NewRouter()

			r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
				GetShortLink(rw, r, s, "http://localhost:8080")
			})

			req := httptest.NewRequest(http.MethodPost, "/", tt.body)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestGetShortLinkAPI(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name string
		body io.Reader
		want want
	}{
		{
			name: "test 1",
			body: strings.NewReader("yandex.ru"),
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name: "test 2",
			body: strings.NewReader(""),
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.New("")

			r := chi.NewRouter()

			r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
				GetShortLink(rw, r, s, "http://localhost:8080")
			})

			req := httptest.NewRequest(http.MethodPost, "/", tt.body)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

package handlers

import (
	"bytes"
	"github.com/IgorPestretsov/yandex_shortener/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

type want struct {
	code        int
	response    string
	contentType string
}

func TestHandlerGenerator_GetFullLinkByID(t *testing.T) {

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		want want
		args args
	}{
		{
			name: "test1_400",
			want: want{
				code:        400,
				response:    `{"status":"ok"}`,
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.Storage{}
			s.Run()
			hg := &HandlerGenerator{
				lastHandler: nil,
				channels:    s.Channels,
			}
			hg.lastHandler = http.HandlerFunc(hg.GetFullLinkByID)
			request := httptest.NewRequest(http.MethodGet, "/?id=XlBz", nil)
			w := httptest.NewRecorder()
			h := hg.GetShortLink(hg.lastHandler)

			h.ServeHTTP(w, request)
			res := w.Result()
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

		})
	}
}

func TestHandlerGenerator_GetShortLink(t *testing.T) {

	tests := []struct {
		name string
		want want
	}{{
		name: "test1_200ok",
		want: want{
			code:        201,
			response:    `{"status":"ok"}`,
			contentType: "application/json",
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.Storage{}
			s.Run()
			hg := &HandlerGenerator{
				lastHandler: nil,
				channels:    s.Channels,
			}
			hg.lastHandler = http.HandlerFunc(hg.GetFullLinkByID)
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("yandex.ru")))
			w := httptest.NewRecorder()
			h := hg.GetShortLink(hg.lastHandler)

			h.ServeHTTP(w, request)
			res := w.Result()
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

		})

	}
}

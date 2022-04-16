package handlers

import (
	"fmt"
	"github.com/IgorPestretsov/yandex_shortener/internal/app"
	"github.com/IgorPestretsov/yandex_shortener/internal/server"
	"github.com/IgorPestretsov/yandex_shortener/internal/storage"
	"io"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func GetFullLinkByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here")
	switch r.Method {

	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "The query parameter is missing", http.StatusBadRequest)
			return
		}

		storage.KeyChannel <- id
		FullLink := <-storage.FullLinkChannel
		if FullLink == "" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", FullLink)
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func GetShortLink(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			b, _ := io.ReadAll(r.Body)
			shortLink := app.GenerateShortLink()
			w.WriteHeader(http.StatusCreated)
			_, err := w.Write([]byte("http://" + server.ServerURL + "/?id=" + shortLink))
			if err != nil {
				return
			}
			storage.LinksPairsChannel <- [2]string{string(b), shortLink}
		default:
			next.ServeHTTP(w, r)
		}
	})
}

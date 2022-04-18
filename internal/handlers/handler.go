package handlers

import (
	"fmt"
	"github.com/IgorPestretsov/yandex_shortener/internal/app"
	"github.com/IgorPestretsov/yandex_shortener/internal/server"
	"github.com/IgorPestretsov/yandex_shortener/internal/storage"
	"io"
	"net/http"
)

type HandlerGenerator struct {
	lastHandler http.Handler
	channels    storage.Channels
}

func (hg *HandlerGenerator) Create(channels storage.Channels) http.Handler {
	hg.lastHandler = http.HandlerFunc(hg.GetFullLinkByID)
	hg.channels = channels
	return Conveyor(hg.lastHandler, hg.GetShortLink)
}

type Middleware func(http.Handler) http.Handler

func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func (hg *HandlerGenerator) GetFullLinkByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	switch r.Method {

	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "The query parameter is missing", http.StatusBadRequest)
			return
		}
		hg.channels.KeyChannel <- id
		FullLink := <-hg.channels.FullLinkChannel
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

func (hg *HandlerGenerator) GetShortLink(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		switch r.Method {
		case http.MethodPost:
			b, _ := io.ReadAll(r.Body)
			shortLink := app.GenerateShortLink()
			fmt.Println(shortLink)
			w.WriteHeader(http.StatusCreated)
			_, err := w.Write([]byte("http://" + server.ServerURL + "/?id=" + shortLink))
			if err != nil {
				return
			}
			hg.channels.LinksPairsChannel <- [2]string{string(b), shortLink}
		default:
			next.ServeHTTP(w, r)
		}
	})
}

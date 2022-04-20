package handlers

import (
	"github.com/IgorPestretsov/yandex_shortener/internal/app"
	"github.com/IgorPestretsov/yandex_shortener/internal/server"
	"github.com/IgorPestretsov/yandex_shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func GetFullLinkByID(w http.ResponseWriter, r *http.Request, s *storage.Storage) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID param is missed", http.StatusBadRequest)
		return
	}
	FullLink := s.LoadLinksPair(id)
	if FullLink == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", FullLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
	r.Body.Close()

}

func GetShortLink(rw http.ResponseWriter, r *http.Request, s *storage.Storage) {
	b, _ := io.ReadAll(r.Body)
	if string(b) == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	shortLink := app.GenerateShortLink()
	s.SaveLinksPair(string(b), shortLink)
	rw.WriteHeader(http.StatusCreated)
	_, err := rw.Write([]byte("http://" + server.ServerURL + "/" + shortLink))
	if err != nil {
		return
	}
}
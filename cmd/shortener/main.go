package main

import (
	"github.com/IgorPestretsov/yandex_shortener/internal/handlers"
	"github.com/IgorPestretsov/yandex_shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	s := storage.New()
	r.Get("/{id}", func(rw http.ResponseWriter, r *http.Request) {
		handlers.GetFullLinkByID(rw, r, s)
	})

	r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
		handlers.GetShortLink(rw, r, s)
	})

	log.Fatal(http.ListenAndServe("localhost:8080", r))

}

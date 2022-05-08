package main

import (
	"fmt"
	"github.com/IgorPestretsov/yandex_shortener/internal/handlers"
	"github.com/IgorPestretsov/yandex_shortener/internal/storage"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:35217"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:35217"`
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	r := chi.NewRouter()
	fmt.Println(cfg)
	s := storage.New()
	r.Get("/{id}", func(rw http.ResponseWriter, r *http.Request) {
		handlers.GetFullLinkByID(rw, r, s)
	})

	r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
		handlers.GetShortLink(rw, r, s, cfg.BaseURL)
	})

	r.Post("/api/shorten", func(rw http.ResponseWriter, r *http.Request) {
		handlers.GetShortLinkAPI(rw, r, s, cfg.BaseURL)
	})
	fmt.Println(cfg.ServerAddress)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))

}

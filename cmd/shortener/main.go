package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"github.com/IgorPestretsov/yandex_shortener/internal/handlers"
	"github.com/IgorPestretsov/yandex_shortener/internal/storage"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	parseFlags(&cfg)

	s := storage.New(cfg.FileStoragePath)
	defer s.Close()

	r := chi.NewRouter()
	r.Use(middleware.Compress(5))
	r.Use(Decompress)
	r.Get("/{id}", func(rw http.ResponseWriter, r *http.Request) {
		handlers.GetFullLinkByID(rw, r, s)
	})

	r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
		handlers.GetShortLink(rw, r, s, cfg.BaseURL)
	})

	r.Post("/api/shorten", func(rw http.ResponseWriter, r *http.Request) {
		handlers.GetShortLinkAPI(rw, r, s, cfg.BaseURL)
	})
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))

}

func parseFlags(config *Config) {
	flag.StringVar(&config.ServerAddress, "a", config.ServerAddress, "Server address to listen on")
	flag.StringVar(&config.BaseURL, "b", config.BaseURL, "Base URL for shortlinks")
	flag.StringVar(&config.FileStoragePath, "f", config.FileStoragePath, "File storage path")
	flag.Parse()
}

func Decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reader io.Reader
		if r.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			reader = gz
			defer gz.Close()
		} else {
			reader = r.Body
		}
		body, err := io.ReadAll(reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(w, r)
	})
}

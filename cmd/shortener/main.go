package main

import (
	"github.com/IgorPestretsov/yandex_shortener/internal/handlers"
	"github.com/IgorPestretsov/yandex_shortener/internal/server"
	"github.com/IgorPestretsov/yandex_shortener/internal/storage"
	"net/http"
)

func main() {

	s := storage.Storage{}
	s.Run()

	hg := handlers.HandlerGenerator{}
	handler := hg.Create(s.Channels)

	http.Handle("/", handler)
	srv := &http.Server{
		Addr: server.ServerURL,
	}
	err := srv.ListenAndServe()
	if err != nil {
		return
	}

}

package main

import (
	"github.com/IgorPestretsov/yandex_shortener/internal/handlers"
	"github.com/IgorPestretsov/yandex_shortener/internal/server"
	"github.com/IgorPestretsov/yandex_shortener/internal/storage"
	"net/http"
)

func main() {
	storage.LinksPairsChannel = make(chan [2]string)
	storage.KeyChannel = make(chan string)
	storage.Storage = make(map[string]string)
	storage.FullLinkChannel = make(chan string)

	GetFullLinkByIDHandler := http.HandlerFunc(handlers.GetFullLinkByID)

	go storage.SaveLinksPair()
	go storage.LoadLinksPair()
	http.Handle("/", handlers.Conveyor(GetFullLinkByIDHandler, handlers.GetShortLink))
	srv := &http.Server{
		Addr: server.ServerURL,
	}
	err := srv.ListenAndServe()
	if err != nil {
		return
	}

}

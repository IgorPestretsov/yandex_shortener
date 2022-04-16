package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

var (
	LinksPairsChannel chan [2]string
	FullLinkChannel   chan string
	KeyChannel        chan string
	Storage           map[string]string
)

const ServerURL = "localhost:8080"

type Middleware func(http.Handler) http.Handler

func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func GenerateShortLink() string {
	const length = 5
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetShortLink(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			b, _ := io.ReadAll(r.Body)
			shortLink := GenerateShortLink()
			w.WriteHeader(http.StatusCreated)
			_, err := w.Write([]byte(ServerURL + "/?id=" + shortLink))
			if err != nil {
				return
			}
			LinksPairsChannel <- [2]string{string(b), shortLink}
		default:
			next.ServeHTTP(w, r)
		}
	})
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

		KeyChannel <- id
		FullLink := <-FullLinkChannel
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

func LoadLinksPair() {
	for {
		key := <-KeyChannel
		FullLink := Storage[key]
		FullLinkChannel <- FullLink

	}
}

func SaveLinksPair() {
	for {
		LinksPair := <-LinksPairsChannel
		Storage[LinksPair[1]] = LinksPair[0]
		fmt.Println(Storage)
	}
}

func main() {
	LinksPairsChannel = make(chan [2]string)
	KeyChannel = make(chan string)
	Storage = make(map[string]string)
	FullLinkChannel = make(chan string)

	GetFullLinkByIDHandler := http.HandlerFunc(GetFullLinkByID)

	go SaveLinksPair()
	go LoadLinksPair()
	http.Handle("/", Conveyor(GetFullLinkByIDHandler, GetShortLink))
	server := &http.Server{
		Addr: ServerURL,
	}
	err := server.ListenAndServe()
	if err != nil {
		return
	}

}

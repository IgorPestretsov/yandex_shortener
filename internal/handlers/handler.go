package handlers

import (
	"encoding/json"
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
func GetShortLinkAPI(rw http.ResponseWriter, r *http.Request, s *storage.Storage, baseURL string) {
	inputData := struct {
		URL string `json:"url"`
	}{}
	GeneratedData := struct {
		Result string `json:"result"`
	}{}
	rawData, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(rawData, &inputData)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	id := app.GenerateShortLink()
	s.SaveLinksPair(inputData.URL, id)
	GeneratedData.Result = baseURL + id

	output, err := json.Marshal(GeneratedData)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	_, err = rw.Write(output)
	if err != nil {
		return
	}
}
func GetShortLinkAPI(rw http.ResponseWriter, r *http.Request, s *storage.Storage) {
	inputData := struct {
		URL string `json:"url"`
	}{}
	GeneratedData := struct {
		Result string `json:"result"`
	}{}
	rawData, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(rawData, &inputData)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(inputData.URL)
	id := app.GenerateShortLink()
	s.SaveLinksPair(inputData.URL, id)
	GeneratedData.Result = "http://" + server.ServerURL + "/" + id

	output, err := json.Marshal(GeneratedData)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(string(output))
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	_, err = rw.Write(output)
	if err != nil {
		return
	}
}

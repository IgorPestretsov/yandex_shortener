package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IgorPestretsov/yandex_shortener/internal/app"
	"github.com/IgorPestretsov/yandex_shortener/internal/middlewares"
	"github.com/IgorPestretsov/yandex_shortener/internal/sqlstorage"
	"github.com/IgorPestretsov/yandex_shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type userRequest struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
type BatchElement struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url,omitempty"`
	ShortURL      string `json:"short_url"`
}

type inputData struct {
	URL string `json:"url"`
}
type generatedData struct {
	Result string `json:"result"`
}

func GetFullLinkByID(w http.ResponseWriter, r *http.Request, s storage.Storage) {
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

func GetShortLink(rw http.ResponseWriter, r *http.Request, s storage.Storage, baseURL string) {
	b, _ := io.ReadAll(r.Body)
	if string(b) == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	uid := r.Context().Value(middlewares.Ctxkey{}).(string)
	shortLink := app.GenerateShortLink()
	existedShortLink, err := s.SaveLinksPair(uid, string(b), shortLink)

	var aee *sqlstorage.AlreadyExistErr
	fmt.Println(err)
	if errors.As(err, &aee) {

		rw.WriteHeader(http.StatusConflict)
		rw.Write([]byte(baseURL + "/" + existedShortLink))
		return
	}
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(baseURL + "/" + shortLink))

}
func GetShortLinkAPI(rw http.ResponseWriter, r *http.Request, s storage.Storage, baseURL string) {

	uid := r.Context().Value(middlewares.Ctxkey{}).(string)
	inData := inputData{}
	genData := generatedData{}
	rawData, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(rawData, &inData)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	id := app.GenerateShortLink()
	genData.Result = baseURL + "/" + id

	existedShortLink, err := s.SaveLinksPair(uid, inData.URL, id)

	rw.Header().Add("Content-Type", "application/json")

	var aee *sqlstorage.AlreadyExistErr
	if errors.As(err, &aee) {

		rw.WriteHeader(http.StatusConflict)
		genData.Result = baseURL + "/" + existedShortLink
	} else {
		rw.WriteHeader(http.StatusCreated)
		genData.Result = baseURL + "/" + id

	}

	output, err := json.Marshal(genData)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = rw.Write(output)
	if err != nil {
		return
	}
}

func GetShortsLinksBatch(rw http.ResponseWriter, r *http.Request, s storage.Storage, baseURL string) {

	uid := r.Context().Value(middlewares.Ctxkey{}).(string)
	var Data []BatchElement
	//genData := generatedData{}
	rawData, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(rawData, &Data)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	for n := range Data {
		Data[n].ShortURL = app.GenerateShortLink()
	}
	for n := range Data {
		s.SaveLinksPair(uid, Data[n].OriginalURL, Data[n].ShortURL)
		Data[n].ShortURL = baseURL + "/" + Data[n].ShortURL
		Data[n].OriginalURL = ""
	}

	output, err := json.Marshal(Data)
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

func GetUserURLs(w http.ResponseWriter, r *http.Request, s storage.Storage, baseURL string) {
	uid := r.Context().Value(middlewares.Ctxkey{}).(string)
	data := s.GetAllUserURLs(uid)
	if len(data) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	var userRequests []*userRequest
	for key, value := range data {
		userRequests = append(userRequests, &userRequest{ShortURL: baseURL + "/" + key, OriginalURL: value})
	}

	output, _ := json.Marshal(userRequests)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err := w.Write(output)
	if err != nil {
		return
	}

}

func PingDB(
	w http.ResponseWriter,
	r *http.Request,
	dsn string,
) {

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

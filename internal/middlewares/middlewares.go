package middlewares

import (
	"bytes"
	"compress/gzip"
	"context"
	"github.com/IgorPestretsov/yandex_shortener/internal/app"
	"io"
	"io/ioutil"
	"net/http"
)

type key string

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

func AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			uid         string
			cookieValue string
		)
		idCookie, err := r.Cookie("uid")
		if err != nil {
			uid, cookieValue = app.GenerateNewUserCookie()
			cookie := http.Cookie{Name: "uid", Value: cookieValue}
			http.SetCookie(w, &cookie)
		} else {
			cookieValue = idCookie.Value
			uid, err = app.GetUserIDfromCookie(cookieValue)
			if err != nil {
				uid, cookieValue = app.GenerateNewUserCookie()
				cookie := http.Cookie{Name: "uid", Value: cookieValue}
				http.SetCookie(w, &cookie)
			}

		}
		var k key = "uuid"
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), k, uid)))
	})

}

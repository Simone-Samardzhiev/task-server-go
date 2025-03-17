package utils

import (
	"log"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request URL:", r.URL.Path)
		log.Println("Request Method:", r.Method)

		next.ServeHTTP(w, r)
	})
}

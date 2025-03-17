package utils

import (
	"log"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request URL: %s \n Request Method: %s\n", r.URL.Path, r.Method)
		next.ServeHTTP(w, r)
	})
}

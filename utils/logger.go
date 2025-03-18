package utils

import (
	"fmt"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request URL: %s, Request Method: %s\n", r.URL.Path, r.Method)
		next.ServeHTTP(w, r)
	})
}

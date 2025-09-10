package middleware

import (
	"log"
	"net/http"
)

// Middleware for logging requests
func Logger(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s endpoint hit", r.Method, r.URL.Path)

		handler(w, r)
	}
}

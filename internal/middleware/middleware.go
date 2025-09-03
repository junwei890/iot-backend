package middleware

import (
	"iot-backend/internal/handlers"
	"iot-backend/pkg/utils"
	"log"
	"net/http"
)

func Logger(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s endpoint hit", r.Method, r.URL.Path)

		handler(w, r)
	}
}

func Authenticator(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusCode, err := utils.ValidateBearerToken(r.Header)
		if err != nil {
			handlers.ErrorWriter(w, statusCode, err)
			return
		}

		handler(w, r)
	}
}

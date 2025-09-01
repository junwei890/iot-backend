package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const authHeaderKey = "Authorization"

func ValidateBearerToken(headers http.Header) (int, error) {
	authHeaderValue := headers.Get(authHeaderKey)
	if authHeaderValue == "" {
		return http.StatusBadRequest, fmt.Errorf("authorization header not present")
	}

	bearerToken := strings.Fields(authHeaderValue)[1]
	if err := godotenv.Load(); err != nil {
		log.Printf("couldn't load environment variables from .env file: %v", err)
	}
	token := os.Getenv("AUTH_TOKEN")
	if bearerToken != token {
		return http.StatusForbidden, fmt.Errorf("bearer token provided is invalid")
	}

	return http.StatusOK, nil
}

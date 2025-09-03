package utils

import (
	"fmt"
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

	authHeaderFields := strings.Fields(authHeaderValue)
	if len(authHeaderFields) != 2 {
		return http.StatusBadRequest, fmt.Errorf("authorization header formatting invalid")
	}
	if authHeaderFields[0] != "Bearer" {
		return http.StatusBadRequest, fmt.Errorf("authorization header formatting invalid")
	}
	bearerToken := authHeaderFields[1]
	godotenv.Load("../.env") // #nosec G104
	token := os.Getenv("AUTH_TOKEN")
	if bearerToken != token {
		return http.StatusForbidden, fmt.Errorf("bearer token provided is invalid")
	}

	return http.StatusOK, nil
}

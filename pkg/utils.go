package pkg

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const issuer = "iot-backend"

// Generates jwt/refresh token with the specified lifetime
func GenerateToken(id, key string, duration time.Duration) (string, error) {
	var token *jwt.Token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Validates token
func ValidateToken(id, key string, r *http.Request) (bool, error) {
	var authHeader string = r.Header.Get("Authorization")
	if authHeader == "" {
		return false, fmt.Errorf("authorization header not provided")
	}
	var headerSlice []string = strings.Fields(authHeader)
	if len(headerSlice) != 2 || headerSlice[0] != "Bearer" {
		return false, fmt.Errorf("invalid authorization header format")
	}
	var tokenString string = headerSlice[1]

	// Only checks hash, other checks must be done post parsing
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(key), nil
	})
	if err != nil {
		return false, err
	}

	// Other checks here
	date, err := token.Claims.GetExpirationTime()
	if err != nil {
		return false, err
	}
	if ok := time.Now().After(date.Time); ok {
		return false, fmt.Errorf("token expired")
	}

	retrievedIssuer, err := token.Claims.GetIssuer()
	if err != nil {
		return false, err
	}
	if retrievedIssuer != issuer {
		return false, fmt.Errorf("invalid token issuer")
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return false, err
	}
	if subject != id {
		return false, fmt.Errorf("invalid token subject")
	}

	return true, nil
}

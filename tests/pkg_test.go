package tests

import (
	"crypto/rand"
	"fmt"
	"iot-backend/pkg"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const authHeaderKey = "Authorization"
const randomInsert = "hamburger"

func TestGenerateAndValidateToken(t *testing.T) {
	buffer := make([]byte, 64)
	if _, err := rand.Read(buffer); err != nil {
		t.Errorf("couldn't setup test: %v", err)
	}
	testKey := string(buffer)
	var testID string = uuid.New().String()

	// Test valid token
	var lifetime time.Duration = 60 * time.Second
	token, _ := pkg.GenerateToken(testID, testKey, lifetime)
	req := &http.Request{
		Header: http.Header{},
	}
	req.Header.Set(authHeaderKey, fmt.Sprintf("Bearer %s", token))
	valid, err := pkg.ValidateToken(testID, testKey, req)
	assert.Nil(t, err)
	assert.True(t, valid)

	// Test expired token
	lifetime = time.Nanosecond
	token, _ = pkg.GenerateToken(testID, testKey, lifetime)
	req = &http.Request{
		Header: http.Header{},
	}
	req.Header.Set(authHeaderKey, fmt.Sprintf("Bearer %s", token))
	valid, err = pkg.ValidateToken(testID, testKey, req)
	assert.NotNil(t, err)
	assert.False(t, valid)

	// Test invalid subject
	lifetime = 60 * time.Second
	token, _ = pkg.GenerateToken(testID, testKey, lifetime)
	req = &http.Request{
		Header: http.Header{},
	}
	req.Header.Set(authHeaderKey, fmt.Sprintf("Bearer %s", token))
	valid, err = pkg.ValidateToken(randomInsert, testKey, req)
	assert.NotNil(t, err)
	assert.False(t, valid)

	// Test invalid key
	lifetime = 60 * time.Second
	token, _ = pkg.GenerateToken(testID, randomInsert, lifetime)
	req = &http.Request{
		Header: http.Header{},
	}
	req.Header.Set(authHeaderKey, fmt.Sprintf("Bearer %s", token))
	valid, err = pkg.ValidateToken(testID, testKey, req)
	assert.NotNil(t, err)
	assert.False(t, valid)

	// Test invalid issuer
	var testToken *jwt.Token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    randomInsert,
		Subject:   randomInsert,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(lifetime)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	tokenString, _ := testToken.SignedString([]byte(testKey))
	req = &http.Request{
		Header: http.Header{},
	}
	req.Header.Set(authHeaderKey, fmt.Sprintf("Bearer %s", tokenString))
	valid, err = pkg.ValidateToken(randomInsert, testKey, req)
	assert.NotNil(t, err)
	assert.False(t, valid)

	// Test missing header
	token, _ = pkg.GenerateToken(testID, testKey, lifetime)
	req = &http.Request{
		Header: http.Header{},
	}
	valid, err = pkg.ValidateToken(testID, testKey, req)
	assert.NotNil(t, err)
	assert.False(t, valid)

	// Test invalid header formatting
	req.Header.Set(authHeaderKey, randomInsert)
	valid, err = pkg.ValidateToken(testID, testKey, req)
	assert.NotNil(t, err)
	assert.False(t, valid)

	// Test invalid header formatting
	req.Header.Set(authHeaderKey, fmt.Sprintf("bearer %s", token))
	valid, err = pkg.ValidateToken(testID, testKey, req)
	assert.NotNil(t, err)
	assert.False(t, valid)
}

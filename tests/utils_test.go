package tests

import (
	"fmt"
	"iot-backend/pkg/utils"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestValidateBearerToken(t *testing.T) {
	// Test: "Authorization" header not present
	header := http.Header{}
	code, err := utils.ValidateBearerToken(header)
	assert.Equal(t, code, 400)
	assert.NotNil(t, err)

	// Test: "Authorization" header present but formatting invalid
	header2 := http.Header{}
	header2.Add("Authorization", "test")
	code, err = utils.ValidateBearerToken(header2)
	assert.Equal(t, code, 400)
	assert.NotNil(t, err)

	// Test: "Authorization" header present but formatting invalid
	header3 := http.Header{}
	header3.Add("Authorization", "Token <token>")
	code, err = utils.ValidateBearerToken(header3)
	assert.Equal(t, code, 400)
	assert.NotNil(t, err)

	// Test: "Authorization" header present, formatting valid but token invalid
	header4 := http.Header{}
	header4.Add("Authorization", "Bearer <token>")
	code, err = utils.ValidateBearerToken(header4)
	assert.Equal(t, code, 403)
	assert.NotNil(t, err)

	// Test: Valid bearer token
	godotenv.Load("../.env")
	authToken := os.Getenv("AUTH_TOKEN")
	header5 := http.Header{}
	header5.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	code, err = utils.ValidateBearerToken(header5)
	assert.Equal(t, code, 200)
	assert.Nil(t, err)
}

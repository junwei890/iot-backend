package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"iot-backend/internal/handlers"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// Only testing handlers I can get wrong logic wise. GetData can be tested with Postman requests
func TestPostData(t *testing.T) {
	godotenv.Load("../.env")
	dbToken := os.Getenv("DB_TOKEN")
	dbURL := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")

	testClient, err := influxdb3.New(influxdb3.ClientConfig{
		Host:     dbURL,
		Token:    dbToken,
		Database: dbName,
	})
	if err != nil {
		t.Errorf("couldn't create client for test database: %v", err)
	}
	defer func(client *influxdb3.Client) {
		if err := client.Close(); err != nil {
			log.Fatal(fmt.Errorf("couldn't gracefully shutdown test client: %v", err))
		}
	}(testClient)

	shared := handlers.Shared{
		DBClient: testClient,
	}

	// Test: Body is nil
	req := httptest.NewRequest(http.MethodPost, "/data", nil)
	writer := httptest.NewRecorder()
	shared.PostData(writer, req)
	res := writer.Result()
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, 400)

	// Test: Body doesn't match expected json structure
	reading := struct {
		Measurement int `json:"measurement"`
	}{
		Measurement: 69,
	}
	readingInBytes, err := json.Marshal(reading)
	if err != nil {
		t.Errorf("couldn't marshal test struct into json: %v", err)
	}
	req2 := httptest.NewRequest(http.MethodPost, "/data", bytes.NewReader(readingInBytes))
	writer2 := httptest.NewRecorder()
	shared.PostData(writer2, req2)
	res2 := writer2.Result()
	defer res2.Body.Close()
	assert.Equal(t, res2.StatusCode, 400)

	// Test: Valid request
	reading2 := struct {
		Measurement string    `json:"measurement" lp:"measurement"`
		Location    string    `json:"location" lp:"tag,location"`
		Temp        float64   `json:"temp" lp:"field,temp"`
		RH          float64   `json:"rh" lp:"field,rh"`
		Radiance    int64     `json:"radiance" lp:"field,radiance"`
		Time        time.Time `json:"timestamp" lp:"timestamp"`
	}{
		Measurement: "measurement",
		Location:    "lab",
		Temp:        30.1,
		RH:          10.5,
		Radiance:    11,
		Time:        time.Now(),
	}
	reading2InBytes, err := json.Marshal(reading2)
	if err != nil {
		t.Errorf("couldn't marshal test struct into json: %v", err)
	}
	req3 := httptest.NewRequest(http.MethodPost, "/data", bytes.NewReader(reading2InBytes))
	writer3 := httptest.NewRecorder()
	shared.PostData(writer3, req3)
	res3 := writer3.Result()
	defer res3.Body.Close()
	assert.Equal(t, res3.StatusCode, 200)
}

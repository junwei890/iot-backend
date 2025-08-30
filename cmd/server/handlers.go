package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
)

const errorHeaderKey = "Content-Type"
const errorHeaderValue = "application/json"

type errorResponse struct {
	Err error `json:"error"`
}

// #nosec G104
// Helper function to handle responses when handlers error
func errorWriter(w http.ResponseWriter, code int, error error) {
	log.Println(error)

	w.Header().Set(errorHeaderKey, errorHeaderValue)
	w.WriteHeader(code)

	resp := errorResponse{
		Err: error,
	}
	respInBytes, _ := json.Marshal(resp)
	w.Write(respInBytes)
}

const healthHeaderKey = "Content-Type"
const healthHeaderValue = "text/html"
const healthResponse = `
<html>
	<head>
		<title>Health Check</title>
	</head>
	<body>
		<h1>Server ready for requests</h1>
	</body>
</html>
`

// Server health check
func (s *shared) healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(healthHeaderKey, healthHeaderValue)
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(healthResponse)); err != nil {
		errorWriter(w, http.StatusInternalServerError, fmt.Errorf("healthz: couldn't write response body: %v", err))
	}
}

type stmReadings struct {
	Measurement string    `json:"measurement" lp:"measurement"`
	Location    string    `json:"location" lp:"tag,location"`
	Temp        float64   `json:"temp" lp:"field,temp"`
	RH          float64   `json:"rh" lp:"field,rh"`
	Radiance    int64     `json:"radiance" lp:"field,radiance"`
	Time        time.Time `json:"timestamp" lp:"timestamp"`
}

// Receiving data from the STM and storing it
func (s *shared) postData(w http.ResponseWriter, r *http.Request) {
	readingsInBytes, err := io.ReadAll(r.Body)
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, fmt.Errorf("postData: couldn't read request: %v", err))
		return
	}
	readings := &stmReadings{}
	if err := json.Unmarshal(readingsInBytes, readings); err != nil {
		errorWriter(w, http.StatusBadRequest, fmt.Errorf("postData: couldn't unmarshal request body: %v", err))
		return
	}

	if err := s.dbClient.WriteData(r.Context(), []any{readings}); err != nil {
		errorWriter(w, http.StatusInternalServerError, fmt.Errorf("postData: couldn't insert data into database: %v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

const getDataQuery = `
SELECT *
FROM readings
WHERE
	time >= now() - interval '5 minutes'
	AND
	location IN ('lab')
`
const getDataHeaderKey = "Content-Type"
const getDataHeaderValue = "application/json"

type dataPoint struct {
	Location string    `json:"location"`
	Temp     float64   `json:"temp"`
	RH       float64   `json:"rh"`
	Radiance int64     `json:"radiance"`
	Time     time.Time `json:"time"`
}

type getDataResp struct {
	Data []dataPoint `json:"data"`
}

// Endpoint for the dashboard
func (s *shared) getData(w http.ResponseWriter, r *http.Request) {
	queriedData, err := s.dbClient.Query(r.Context(), getDataQuery)
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, fmt.Errorf("getData: couldn't query data from database: %v", err))
		return
	}

	resp := getDataResp{}
	for queriedData.Next() {
		value := queriedData.Value()
		dataPoint := dataPoint{
			Location: value["location"].(string),
			Temp:     value["temp"].(float64),
			RH:       value["rh"].(float64),
			Radiance: value["radiance"].(int64),
			Time:     value["time"].(arrow.Timestamp).ToTime(arrow.Nanosecond),
		}
		resp.Data = append(resp.Data, dataPoint)
	}

	w.Header().Set(getDataHeaderKey, getDataHeaderValue)
	w.WriteHeader(http.StatusOK)

	respInBytes, err := json.Marshal(resp)
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, fmt.Errorf("getData: couldn't marshal response: %v", err))
		return
	}
	if _, err := w.Write(respInBytes); err != nil {
		errorWriter(w, http.StatusInternalServerError, fmt.Errorf("getData: couldn't write response body: %v", err))
	}
}

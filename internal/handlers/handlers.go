package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/apache/arrow/go/v15/arrow"
)

type Shared struct {
	DBClient *influxdb3.Client
}

const errorHeaderKey = "Content-Type"
const errorHeaderValue = "application/json"

type errorResponse struct {
	Err string `json:"error"`
}

// #nosec G104
// Helper function to handle responses when handlers error
func ErrorWriter(w http.ResponseWriter, code int, error error) {
	log.Println(error)

	w.Header().Set(errorHeaderKey, errorHeaderValue)
	w.WriteHeader(code)

	// error is not a json type, needs to be converted to a string
	resp := errorResponse{
		Err: error.Error(),
	}
	respInBytes, _ := json.Marshal(resp)
	w.Write(respInBytes)
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
func (s *Shared) PostData(w http.ResponseWriter, r *http.Request) {
	readingsInBytes, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorWriter(w, http.StatusInternalServerError, fmt.Errorf("%s %s: couldn't read request: %v", r.Method, r.URL.Path, err))
		return
	}
	readings := &stmReadings{}
	if err := json.Unmarshal(readingsInBytes, readings); err != nil {
		ErrorWriter(w, http.StatusBadRequest, fmt.Errorf("%s %s: couldn't unmarshal request body: %v", r.Method, r.URL.Path, err))
		return
	}

	if err := s.DBClient.WriteData(r.Context(), []any{readings}); err != nil {
		ErrorWriter(w, http.StatusInternalServerError, fmt.Errorf("%s %s: couldn't insert data into database: %v", r.Method, r.URL.Path, err))
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
func (s *Shared) GetData(w http.ResponseWriter, r *http.Request) {
	queriedData, err := s.DBClient.Query(r.Context(), getDataQuery)
	if err != nil {
		ErrorWriter(w, http.StatusInternalServerError, fmt.Errorf("%s %s: couldn't query data from database: %v", r.Method, r.URL.Path, err))
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
		ErrorWriter(w, http.StatusInternalServerError, fmt.Errorf("%s %s: couldn't marshal response: %v", r.Method, r.URL.Path, err))
		return
	}
	if _, err := w.Write(respInBytes); err != nil {
		ErrorWriter(w, http.StatusInternalServerError, fmt.Errorf("%s %s: couldn't write response body: %v", r.Method, r.URL.Path, err))
	}
}

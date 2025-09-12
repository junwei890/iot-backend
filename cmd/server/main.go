package main

import (
	"fmt"
	"iot-backend/internal/handlers"
	"iot-backend/internal/middleware"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/joho/godotenv"
)

func main() {
	// Creating DB client
	godotenv.Load() // #nosec G104
	var dbToken string = os.Getenv("DB_TOKEN")
	var dbURL string = os.Getenv("DB_URL")
	var dbName string = os.Getenv("DB_NAME")

	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:     dbURL,
		Token:    dbToken,
		Database: dbName,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("couldn't create client for database: %v", err))
	}
	defer func(client *influxdb3.Client) {
		if err := client.Close(); err != nil {
			log.Fatal(fmt.Errorf("couldn't gracefully shutdown client: %v", err))
		}
	}(client)

	shared := handlers.Shared{
		DBClient: client,
	}

	// Registering endpoints
	var mux *http.ServeMux = http.NewServeMux()
	mux.HandleFunc("POST /data", middleware.Logger(shared.PostData))
	mux.HandleFunc("GET /data", middleware.Logger(shared.GetData))

	// Starting server
	var port string = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := http.Server{
		Addr:              fmt.Sprintf("0.0.0.0:%s", port),
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		MaxHeaderBytes:    8192,
	}
	log.Printf("started server on port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(fmt.Errorf("couldn't start server: %v", err))
	}
}

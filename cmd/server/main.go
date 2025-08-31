package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/joho/godotenv"
)

type shared struct {
	dbClient *influxdb3.Client
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("couldn't load environment variables from .env file: %v", err)
	}

	dbToken := os.Getenv("DB_TOKEN")
	dbURL := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")

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

	shared := shared{
		dbClient: client,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", shared.healthz)
	mux.HandleFunc("POST /data", shared.postData)
	mux.HandleFunc("GET /data", shared.getData)

	port := os.Getenv("PORT")
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

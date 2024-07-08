// main.go

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/joshua468/weather-location-api-master/handler"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port if PORT environment variable is not set
	}

	http.HandleFunc("/api/hello", handler.HelloHandler)

	log.Printf("Server listening on port %s", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

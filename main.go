package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/tomasen/realip"
)

type Response struct {
	ClientIP string `json:"client_ip"`
	Location string `json:"location"`
	Greeting string `json:"greeting"`
}

type IPInfo struct {
	City string `json:"city"`
}

type WeatherInfo struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func getLocation(ip string) (string, error) {
	url := "https://ipinfo.io/" + ip + "/json"
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch geolocation data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var info IPInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return "", fmt.Errorf("failed to decode geolocation data: %v", err)
	}

	if info.City == "" {
		return "Unknown", nil
	}

	return info.City, nil
}

func getWeather(city string, apiKey string) (float64, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s", city, apiKey)
	log.Printf("Fetching weather data from URL: %s", url) // Log the request URL

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch weather data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var weather WeatherInfo
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		return 0, fmt.Errorf("failed to decode weather data: %v", err)
	}

	return weather.Main.Temp, nil
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	visitorName := r.URL.Query().Get("visitor_name")
	if visitorName == "" {
		visitorName = "Guest"
	}

	clientIP := realip.FromRequest(r)
	if clientIP == "::1" || clientIP == "127.0.0.1" {
		clientIP = "8.8.8.8" // Default to a public IP for local testing
	}

	location, err := getLocation(clientIP)
	if err != nil {
		log.Printf("Error getting location: %v", err)
		location = "Unknown"
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable is not set")
	}

	temperature, err := getWeather(location, apiKey)
	if err != nil {
		log.Printf("Error getting weather: %v", err)
		temperature = 11.0 // Default temperature in case of error
	}

	response := Response{
		ClientIP: clientIP,
		Location: location,
		Greeting: fmt.Sprintf("Hello, %s! The temperature is %.1f degrees Celsius in %s", visitorName, temperature, location),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

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

	http.HandleFunc("/api/hello", HelloHandler)

	log.Printf("Server listening on port %s", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

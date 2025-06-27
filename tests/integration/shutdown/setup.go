package shutdown

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Set up test environment variables
	envVars := map[string]string{
		"DB_HOST":                "localhost",
		"DB_PORT":                "5432",
		"DB_USER":                "test",
		"DB_PASSWORD":            "test",
		"DB_NAME":                "test",
		"SENDGRID_API_KEY":       "test_key",
		"SENDER_EMAIL":           "test@example.com",
		"SENDER_NAME":            "Test Sender",
		"OPENWEATHERMAP_API_KEY": "test_key",
		"GEOCODING_API_URL":      "http://localhost:8080",
		"OPENWEATHERMAP_API_URL": "http://localhost:8080",
		"WEATHER_API_KEY":        "test_key",
		"WEATHER_API_URL":        "http://localhost:8080",
	}

	for key, value := range envVars {
		if err := os.Setenv(key, value); err != nil {
			log.Printf("Failed to set environment variable %s: %v", key, err)
		}
	}

	// Run tests
	code := m.Run()

	// Clean up
	os.Exit(code)
}

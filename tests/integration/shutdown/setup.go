package shutdown

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Set up test environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "test")
	os.Setenv("DB_PASSWORD", "test")
	os.Setenv("DB_NAME", "test")
	os.Setenv("SENDGRID_API_KEY", "test_key")
	os.Setenv("SENDER_EMAIL", "test@example.com")
	os.Setenv("SENDER_NAME", "Test Sender")
	os.Setenv("OPENWEATHERMAP_API_KEY", "test_key")
	os.Setenv("GEOCODING_API_URL", "http://localhost:8080")
	os.Setenv("OPENWEATHERMAP_API_URL", "http://localhost:8080")
	os.Setenv("WEATHER_API_KEY", "test_key")
	os.Setenv("WEATHER_API_URL", "http://localhost:8080")

	// Run tests
	code := m.Run()

	// Clean up
	os.Exit(code)
}

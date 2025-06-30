package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	err := godotenv.Load() // Load .env file (if exists)
	if err != nil {
		// Log, but don't fatal, as environment variables might be set directly
		fmt.Printf("No .env file found or error loading it: %v. Using environment variables directly.\n", err)
	}

	return &Config{
		AppPort:    getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", ""),
		DBPort:     getEnv("DB_PORT", ""),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", ""),
	}
}

// getEnv gets environment variable or uses fallback
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

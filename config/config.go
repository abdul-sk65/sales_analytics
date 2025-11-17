package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	MongoURI       string
	DatabaseName   string
	CSVFilePath    string
	Port           string
	WorkerPoolSize int
}

// Load reads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	workerPoolSize := 10
	if wps := os.Getenv("WORKER_POOL_SIZE"); wps != "" {
		if parsed, err := strconv.Atoi(wps); err == nil {
			workerPoolSize = parsed
		}
	}

	return &Config{
		MongoURI:       getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DatabaseName:   getEnv("DATABASE_NAME", "sales_analytics"),
		CSVFilePath:    getEnv("CSV_FILE_PATH", "./data/sales_data.csv"),
		Port:           getEnv("PORT", "8080"),
		WorkerPoolSize: workerPoolSize,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

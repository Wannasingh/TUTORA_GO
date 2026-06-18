package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DBConnString string
	DBSchema     string
}

func LoadConfig() *Config {
	// Load environment variables from .env if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading configurations from system environment variables")
	}

	port := getEnv("PORT", "8080")
	dbConn := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/tutora?sslmode=disable")
	dbSchema := getEnv("DATABASE_SCHEMA", "tutora_app")

	return &Config{
		Port:         port,
		DBConnString: dbConn,
		DBSchema:     dbSchema,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

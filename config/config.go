package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT           string
	DATABASE_URL   string
	HOST           string
	USER           string
	PASSWORD       string
	DB_NAME        string
	DB_PORT        string
	SSLMODE        string
	CORS_ORIGIN    string
	ADMIN_EMAIL    string
	ADMIN_NAME     string
	ADMIN_PASSWORD string
}

var Envs = initConfig()

func initConfig() Config {
	_ = godotenv.Load()

	return Config{
		PORT:           getEnv("PORT", getEnv("HTTP_PORT", "8000")),
		DATABASE_URL:   getEnv("DATABASE_URL", ""),
		HOST:           getEnv("HOST", "localhost"),
		USER:           getEnv("USER", "postgres"),
		PASSWORD:       getEnv("PASSWORD", "postgres"),
		DB_NAME:        getEnv("DB_NAME", "conelli_admin"),
		DB_PORT:        getEnv("DB_PORT", "5432"),
		SSLMODE:        getEnv("SSLMODE", "disable"),
		CORS_ORIGIN:    getEnv("CORS_ORIGIN", "http://localhost:5173"),
		ADMIN_EMAIL:    getEnv("ADMIN_EMAIL", "admin@conelliengineering.com"),
		ADMIN_NAME:     getEnv("ADMIN_NAME", "Conelli Admin"),
		ADMIN_PASSWORD: getEnv("ADMIN_PASSWORD", "change-me"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

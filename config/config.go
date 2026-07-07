package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT              string
	APP_ENV           string
	DATABASE_URL      string
	HOST              string
	DB_USER           string
	PASSWORD          string
	DB_NAME           string
	DB_PORT           string
	SSLMODE           string
	CORS_ORIGIN       string
	ADMIN_EMAIL       string
	ADMIN_NAME        string
	ADMIN_PASSWORD    string
	AWS_REGION        string
	AWS_S3_BUCKET     string
	AWS_S3_PREFIX     string
	AWS_S3_PUBLIC_URL string
	AWS_S3_ENDPOINT   string
}

func (c Config) IsProduction() bool {
	return strings.EqualFold(c.APP_ENV, "production")
}

func (c Config) Validate() error {
	if !c.IsProduction() {
		return nil
	}

	required := map[string]string{
		"PORT":           c.PORT,
		"DATABASE_URL":   c.DATABASE_URL,
		"CORS_ORIGIN":    c.CORS_ORIGIN,
		"ADMIN_EMAIL":    c.ADMIN_EMAIL,
		"ADMIN_PASSWORD": c.ADMIN_PASSWORD,
		"AWS_REGION":     c.AWS_REGION,
		"AWS_S3_BUCKET":  c.AWS_S3_BUCKET,
	}

	missing := make([]string, 0)
	for key, value := range required {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, key)
		}
	}

	if c.ADMIN_PASSWORD == "change-me" {
		missing = append(missing, "ADMIN_PASSWORD must not use the default development value")
	}

	if len(missing) > 0 {
		return fmt.Errorf("invalid production configuration: %s", strings.Join(missing, ", "))
	}

	return nil
}

var Envs = initConfig()

func initConfig() Config {
	_ = godotenv.Load()

	return Config{
		PORT:              getEnv("PORT", getEnv("HTTP_PORT", "8000")),
		APP_ENV:           getEnv("APP_ENV", "development"),
		DATABASE_URL:      getEnv("DATABASE_URL", ""),
		HOST:              getEnv("HOST", "localhost"),
		DB_USER:           getEnv("DB_USER", "myke"),
		PASSWORD:          getEnv("PASSWORD", "password123"),
		DB_NAME:           getEnv("DB_NAME", "conelli_admin"),
		DB_PORT:           getEnv("DB_PORT", "55432"),
		SSLMODE:           getEnv("SSLMODE", "disable"),
		CORS_ORIGIN:       getEnv("CORS_ORIGIN", "http://localhost:5173,http://localhost:5174"),
		ADMIN_EMAIL:       getEnv("ADMIN_EMAIL", "admin@conelliengineering.com"),
		ADMIN_NAME:        getEnv("ADMIN_NAME", "Conelli Admin"),
		ADMIN_PASSWORD:    getEnv("ADMIN_PASSWORD", "change-me"),
		AWS_REGION:        getEnv("AWS_REGION", "eu-west-1"),
		AWS_S3_BUCKET:     getEnv("AWS_S3_BUCKET", ""),
		AWS_S3_PREFIX:     getEnv("AWS_S3_PREFIX", "conelli"),
		AWS_S3_PUBLIC_URL: getEnv("AWS_S3_PUBLIC_URL", ""),
		AWS_S3_ENDPOINT:   getEnv("AWS_S3_ENDPOINT", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"

	"github.com/conelli/admin-backend/config"
	"github.com/conelli/admin-backend/db"
	"github.com/conelli/admin-backend/db/migrations"
)

func main() {
	_ = godotenv.Load()

	direction := flag.String("direction", "up", "migration direction: up or down")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	client, err := db.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer client.Close()

	sqlDB, err := client.SQLDB()
	if err != nil {
		log.Fatalf("failed to access sql database: %v", err)
	}

	if err := migrations.Run(sqlDB, *direction); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}

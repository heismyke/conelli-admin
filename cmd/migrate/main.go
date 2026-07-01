package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/conelli/admin-backend/internal/store"
)

func main() {
	targetNumber := 0
	if len(os.Args) > 2 {
		number, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("invalid migration target: %v", err)
		}
		targetNumber = number
	}

	db, err := store.NewPool()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to access sql database: %v", err)
	}

	if err := store.Migrate(context.Background(), sqlDB, targetNumber); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}

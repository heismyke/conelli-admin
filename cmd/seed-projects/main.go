package main

import (
	"context"
	"log"

	"github.com/conelli/admin-backend/internal/store"
	"github.com/conelli/admin-backend/internal/store/repo"
)

func main() {
	db, err := store.NewPool()
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	adminRepo := repo.NewRepo(store.NewPostgresDBStore(db))
	if err := adminRepo.Init(); err != nil {
		log.Fatalf("failed to initialize admin store: %v", err)
	}

	if err := adminRepo.SeedProjects(context.Background()); err != nil {
		log.Fatalf("failed to seed projects: %v", err)
	}

	log.Println("seeded admin projects from bundled project catalog")
}

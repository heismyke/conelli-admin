package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/conelli/admin-backend/api"
	"github.com/conelli/admin-backend/config"
	"github.com/conelli/admin-backend/db"
	"github.com/conelli/admin-backend/db/repo"
	"github.com/conelli/admin-backend/internal/router"
	"github.com/conelli/admin-backend/logic"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables...")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if cfg.IsDev {
		log.SetLevel(log.DebugLevel)
		log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	} else {
		log.SetLevel(log.InfoLevel)
		log.SetFormatter(&log.JSONFormatter{})
	}

	log.WithFields(log.Fields{
		"port": cfg.HttpPort,
		"dev":  cfg.IsDev,
	}).Info("Starting Conelli Admin Backend API")

	if err := StartServer(cfg); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func StartServer(cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbClient, err := db.New(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}
	defer dbClient.Close()

	if err := dbClient.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	log.Info("Database connected successfully")

	repository := repo.New(dbClient)
	businessLogic := logic.New(cfg, repository)

	healthAPI := api.NewHealth(businessLogic)
	adminAPI := api.NewAdmin(businessLogic)

	r := router.NewRouter(cfg, businessLogic)
	r.RegisterAPI(healthAPI)
	r.RegisterAPI(adminAPI)

	log.Info("API routes registered")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		log.Infof("Received signal: %s. Shutting down gracefully...", sig)
		cancel()
	}()

	log.Infof("Starting HTTP server on port %d", cfg.HttpPort)
	if err := r.Start(ctx); err != nil {
		return fmt.Errorf("server start error: %w", err)
	}

	return nil
}

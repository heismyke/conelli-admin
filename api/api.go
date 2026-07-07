package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/conelli/admin-backend/config"
	"github.com/conelli/admin-backend/internal/handlers/admin"
	"github.com/conelli/admin-backend/internal/handlers/auth"
	"github.com/conelli/admin-backend/internal/handlers/health"
	"github.com/conelli/admin-backend/internal/handlers/realtime"
	"github.com/conelli/admin-backend/internal/store"
	"github.com/conelli/admin-backend/internal/store/repo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Api struct {
	addr   string
	router *gin.Engine
}

func NewApi(addr string) (*Api, error) {
	if config.Envs.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := store.NewPool()
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	postgresStore := store.NewPostgresDBStore(db)
	adminRepo := repo.NewRepo(postgresStore)

	if err := adminRepo.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize admin store: %w", err)
	}

	adminHandler := admin.NewAdminHandler(adminRepo)
	authHandler := auth.NewAuthHandler(adminRepo)
	healthHandler := health.NewHealthHandler()
	realtimeHandler := realtime.NewRealtimeHandler()

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	healthHandler.RegisterHealthHandler(router.Group("/health"))
	authHandler.RegisterAuthHandler(router.Group("/auth"))
	adminHandler.RegisterAdminHandler(router.Group("/admin"))
	realtimeHandler.RegisterRealtimeHandler(router.Group("/realtime"))

	return &Api{
		addr:   addr,
		router: router,
	}, nil
}

func (a *Api) Run() error {
	server := &http.Server{
		Addr:              a.addr,
		Handler:           a.router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-quit:
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	}
}

func corsOrigins() []string {
	origins := strings.Split(config.Envs.CORS_ORIGIN, ",")
	allowed := make([]string, 0, len(origins))
	for _, origin := range origins {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			allowed = append(allowed, origin)
		}
	}

	return allowed
}

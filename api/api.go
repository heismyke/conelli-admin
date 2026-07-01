package api

import (
	"fmt"

	"github.com/conelli/admin-backend/config"
	"github.com/conelli/admin-backend/internal/handlers/admin"
	"github.com/conelli/admin-backend/internal/handlers/auth"
	"github.com/conelli/admin-backend/internal/handlers/health"
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

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.Envs.CORS_ORIGIN},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	healthHandler.RegisterHealthHandler(router.Group("/health"))
	authHandler.RegisterAuthHandler(router.Group("/auth"))
	adminHandler.RegisterAdminHandler(router.Group("/admin"))

	return &Api{
		addr:   addr,
		router: router,
	}, nil
}

func (a *Api) Run() error {
	return a.router.Run(a.addr)
}

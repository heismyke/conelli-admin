package router

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/conelli/admin-backend/config"
	"github.com/conelli/admin-backend/logic"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type API interface {
	RegisterRoutes(engine *gin.Engine)
}

type Router struct {
	cfg    *config.Config
	logic  *logic.CoreLogic
	engine *gin.Engine
	server *http.Server
}

func NewRouter(cfg *config.Config, businessLogic *logic.CoreLogic) *Router {
	if !cfg.IsDev {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.CorsOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return &Router{
		cfg:    cfg,
		logic:  businessLogic,
		engine: engine,
	}
}

func (r *Router) RegisterAPI(api API) {
	api.RegisterRoutes(r.engine)
}

func (r *Router) Start(ctx context.Context) error {
	r.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", r.cfg.HttpPort),
		Handler:           r.engine,
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- r.server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return r.server.Shutdown(shutdownCtx)
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}

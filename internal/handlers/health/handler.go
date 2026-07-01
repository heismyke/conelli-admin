package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHealthHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterHealthHandler(r *gin.RouterGroup) {
	r.GET("", h.GetHealthStatus)
	r.GET("/", h.GetHealthStatus)
	r.GET("/ping", h.Ping)
}

func (h *Handler) GetHealthStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"version": "1.0.0",
		"service": "Conelli Admin API",
	})
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

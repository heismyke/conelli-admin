package auth

import (
	"net/http"

	"github.com/conelli/admin-backend/internal/store/dao"
	"github.com/conelli/admin-backend/internal/store/repo"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	authStore *repo.Repo
}

func NewAuthHandler(s *repo.Repo) *Handler {
	return &Handler{
		authStore: s,
	}
}

func (h *Handler) RegisterAuthHandler(r *gin.RouterGroup) {
	r.POST("/login", h.Login)
	r.POST("/investor/login", h.InvestorLogin)
}

func (h *Handler) Login(c *gin.Context) {
	var payload dao.LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authStore.Login(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) InvestorLogin(c *gin.Context) {
	var payload dao.LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authStore.InvestorLogin(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

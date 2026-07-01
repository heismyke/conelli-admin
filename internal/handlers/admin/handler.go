package admin

import (
	"net/http"

	"github.com/conelli/admin-backend/internal/store/dao"
	"github.com/conelli/admin-backend/internal/store/repo"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	adminStore *repo.Repo
}

func NewAdminHandler(s *repo.Repo) *Handler {
	return &Handler{
		adminStore: s,
	}
}

func (h *Handler) RegisterAdminHandler(r *gin.RouterGroup) {
	r.GET("", h.GetSummary)
	r.GET("/", h.GetSummary)
	r.GET("/data", h.GetData)
	r.PUT("/data", h.SaveData)
	r.POST("/login", h.Login)
}

func (h *Handler) GetSummary(c *gin.Context) {
	c.JSON(http.StatusOK, h.adminStore.AdminSummary(c.Request.Context()))
}

func (h *Handler) GetData(c *gin.Context) {
	data, err := h.adminStore.AdminData(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) SaveData(c *gin.Context) {
	var payload dao.AdminData
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(payload.Users) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one staff user is required"})
		return
	}

	if err := h.adminStore.SaveAdminData(c.Request.Context(), payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data, err := h.adminStore.AdminData(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) Login(c *gin.Context) {
	var payload dao.LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.adminStore.Login(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

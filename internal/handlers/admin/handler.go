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
	r.POST("/uploads/presign", h.PresignUpload)

	r.GET("/properties", h.ListProperties)
	r.POST("/properties", h.CreateProperty)
	r.GET("/properties/:id", h.GetProperty)
	r.PUT("/properties/:id", h.UpdateProperty)
	r.DELETE("/properties/:id", h.DeleteProperty)

	r.GET("/investors", h.ListInvestors)
	r.POST("/investors", h.CreateInvestor)
	r.GET("/investors/:id", h.GetInvestor)
	r.PUT("/investors/:id", h.UpdateInvestor)
	r.DELETE("/investors/:id", h.DeleteInvestor)
	r.PUT("/investors/:id/properties/:propertyId", h.SetInvestorProperty)
	r.DELETE("/investors/:id/properties/:propertyId", h.UnsetInvestorProperty)

	r.GET("/users", h.ListUsers)
	r.POST("/users", h.CreateUser)
	r.GET("/users/:id", h.GetUser)
	r.PUT("/users/:id", h.UpdateUser)
	r.DELETE("/users/:id", h.DeleteUser)

	r.POST("/updates", h.CreateUpdate)
	r.PUT("/updates/:id", h.UpdateUpdate)
	r.DELETE("/updates/:id", h.DeleteUpdate)

	r.POST("/milestones", h.CreateMilestone)
	r.PUT("/milestones/:id", h.UpdateMilestone)
	r.DELETE("/milestones/:id", h.DeleteMilestone)

	r.POST("/documents", h.CreateDocument)
	r.PUT("/documents/:id", h.UpdateDocument)
	r.DELETE("/documents/:id", h.DeleteDocument)
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
	sanitizeInvestors(&data)

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
	current, err := h.adminStore.AdminData(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	preserveInvestorPasswords(&payload, current)

	if err := h.adminStore.SaveAdminData(c.Request.Context(), payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data, err := h.adminStore.AdminData(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	sanitizeInvestors(&data)

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

func sanitizeInvestors(data *dao.AdminData) {
	for i := range data.Investors {
		data.Investors[i].PasswordHash = ""
	}
}

func preserveInvestorPasswords(next *dao.AdminData, current dao.AdminData) {
	passwords := make(map[string]string, len(current.Investors))
	for _, investor := range current.Investors {
		passwords[investor.ID] = investor.PasswordHash
	}
	for i := range next.Investors {
		if next.Investors[i].PasswordHash != "" {
			continue
		}
		if passwordHash := passwords[next.Investors[i].ID]; passwordHash != "" {
			next.Investors[i].PasswordHash = passwordHash
		}
	}
}

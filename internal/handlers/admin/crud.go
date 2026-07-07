package admin

import (
	"net/http"
	"strings"
	"time"

	"github.com/conelli/admin-backend/internal/store/dao"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func newID(prefix string) string {
	return prefix + "_" + strings.ReplaceAll(uuid.NewString(), "-", "")[:12]
}

func nowISO() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

func today() string {
	return time.Now().UTC().Format("2006-01-02")
}

func badRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": message})
}

func notFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{"error": message})
}

func serverError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func (h *Handler) loadData(c *gin.Context) (dao.AdminData, bool) {
	data, err := h.adminStore.AdminData(c.Request.Context())
	if err != nil {
		serverError(c, err)
		return dao.AdminData{}, false
	}
	return data, true
}

func (h *Handler) saveData(c *gin.Context, data dao.AdminData) bool {
	if err := h.adminStore.SaveAdminData(c.Request.Context(), data); err != nil {
		serverError(c, err)
		return false
	}
	return true
}

func findProperty(data dao.AdminData, id string) (*dao.Property, int) {
	for i := range data.Properties {
		if data.Properties[i].ID == id {
			return &data.Properties[i], i
		}
	}
	return nil, -1
}

func findInvestor(data dao.AdminData, id string) (*dao.Investor, int) {
	for i := range data.Investors {
		if data.Investors[i].ID == id {
			return &data.Investors[i], i
		}
	}
	return nil, -1
}

func findUser(data dao.AdminData, id string) (*dao.AdminUser, int) {
	for i := range data.Users {
		if data.Users[i].ID == id {
			return &data.Users[i], i
		}
	}
	return nil, -1
}

func touchProperty(data *dao.AdminData, propertyID string) {
	property, _ := findProperty(*data, propertyID)
	if property != nil {
		property.UpdatedAt = nowISO()
	}
}

func (h *Handler) ListProperties(c *gin.Context) {
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, data.Properties)
}

func (h *Handler) GetProperty(c *gin.Context) {
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	property, _ := findProperty(data, c.Param("id"))
	if property == nil {
		notFound(c, "property not found")
		return
	}
	c.JSON(http.StatusOK, property)
}

func (h *Handler) CreateProperty(c *gin.Context) {
	var payload dao.Property
	if err := c.ShouldBindJSON(&payload); err != nil {
		badRequest(c, err.Error())
		return
	}
	if strings.TrimSpace(payload.Title) == "" || strings.TrimSpace(payload.Location) == "" {
		badRequest(c, "title and location are required")
		return
	}
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	if payload.ID == "" {
		payload.ID = newID("prop")
	}
	if payload.Status == "" {
		payload.Status = "Under Construction"
	}
	if payload.Category == "" {
		payload.Category = "Real Estate Development"
	}
	if payload.EstCompletionDate == "" {
		payload.EstCompletionDate = today()
	}
	if payload.CoverImageURL == "" {
		payload.CoverImageURL = "/assets/home.png"
	}
	if payload.Year == "" {
		payload.Year = time.Now().UTC().Format("2006")
	}
	if len(payload.Tags) == 0 && payload.Category != "" {
		payload.Tags = []string{payload.Category}
	}
	if len(payload.GalleryImages) == 0 && payload.CoverImageURL != "" {
		payload.GalleryImages = []string{payload.CoverImageURL}
	}
	payload.CreatedAt = nowISO()
	payload.UpdatedAt = payload.CreatedAt
	data.Properties = append([]dao.Property{payload}, data.Properties...)
	if !h.saveData(c, data) {
		return
	}
	c.JSON(http.StatusCreated, payload)
}

func (h *Handler) UpdateProperty(c *gin.Context) {
	var payload dao.Property
	if err := c.ShouldBindJSON(&payload); err != nil {
		badRequest(c, err.Error())
		return
	}
	if strings.TrimSpace(payload.Title) == "" || strings.TrimSpace(payload.Location) == "" {
		badRequest(c, "title and location are required")
		return
	}
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	property, _ := findProperty(data, c.Param("id"))
	if property == nil {
		notFound(c, "property not found")
		return
	}
	payload.ID = property.ID
	payload.CreatedAt = property.CreatedAt
	payload.UpdatedAt = nowISO()
	*property = payload
	if !h.saveData(c, data) {
		return
	}
	c.JSON(http.StatusOK, property)
}

func (h *Handler) DeleteProperty(c *gin.Context) {
	id := c.Param("id")
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	_, index := findProperty(data, id)
	if index == -1 {
		notFound(c, "property not found")
		return
	}
	data.Properties = append(data.Properties[:index], data.Properties[index+1:]...)
	data.InvestorProperties = filter(data.InvestorProperties, func(item dao.InvestorProperty) bool { return item.PropertyID != id })
	data.Updates = filter(data.Updates, func(item dao.Update) bool { return item.PropertyID != id })
	data.Milestones = filter(data.Milestones, func(item dao.Milestone) bool { return item.PropertyID != id })
	data.Materials = filter(data.Materials, func(item dao.Material) bool { return item.PropertyID != id })
	data.Documents = filter(data.Documents, func(item dao.Document) bool { return item.PropertyID == nil || *item.PropertyID != id })
	if !h.saveData(c, data) {
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ListInvestors(c *gin.Context) {
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, data.Investors)
}

func (h *Handler) GetInvestor(c *gin.Context) {
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	investor, _ := findInvestor(data, c.Param("id"))
	if investor == nil {
		notFound(c, "investor not found")
		return
	}
	c.JSON(http.StatusOK, investor)
}

func (h *Handler) CreateInvestor(c *gin.Context) {
	var payload dao.Investor
	if err := c.ShouldBindJSON(&payload); err != nil {
		badRequest(c, err.Error())
		return
	}
	if strings.TrimSpace(payload.Name) == "" || strings.TrimSpace(payload.Email) == "" {
		badRequest(c, "name and email are required")
		return
	}
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	payload.ID = newID("inv")
	if payload.MemberSince == "" {
		payload.MemberSince = today()
	}
	if payload.Status == "" {
		payload.Status = "active"
	}
	data.Investors = append([]dao.Investor{payload}, data.Investors...)
	if !h.saveData(c, data) {
		return
	}
	c.JSON(http.StatusCreated, payload)
}

func (h *Handler) UpdateInvestor(c *gin.Context) {
	var payload dao.Investor
	if err := c.ShouldBindJSON(&payload); err != nil {
		badRequest(c, err.Error())
		return
	}
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	investor, _ := findInvestor(data, c.Param("id"))
	if investor == nil {
		notFound(c, "investor not found")
		return
	}
	payload.ID = investor.ID
	*investor = payload
	if !h.saveData(c, data) {
		return
	}
	c.JSON(http.StatusOK, investor)
}

func (h *Handler) DeleteInvestor(c *gin.Context) {
	id := c.Param("id")
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	_, index := findInvestor(data, id)
	if index == -1 {
		notFound(c, "investor not found")
		return
	}
	data.Investors = append(data.Investors[:index], data.Investors[index+1:]...)
	data.InvestorProperties = filter(data.InvestorProperties, func(item dao.InvestorProperty) bool { return item.InvestorID != id })
	data.Documents = filter(data.Documents, func(item dao.Document) bool { return item.InvestorID == nil || *item.InvestorID != id })
	if !h.saveData(c, data) {
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) SetInvestorProperty(c *gin.Context) {
	investorID := c.Param("id")
	propertyID := c.Param("propertyId")
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	if investor, _ := findInvestor(data, investorID); investor == nil {
		notFound(c, "investor not found")
		return
	}
	if property, _ := findProperty(data, propertyID); property == nil {
		notFound(c, "property not found")
		return
	}
	for _, assignment := range data.InvestorProperties {
		if assignment.InvestorID == investorID && assignment.PropertyID == propertyID {
			c.JSON(http.StatusOK, assignment)
			return
		}
	}
	assignment := dao.InvestorProperty{InvestorID: investorID, PropertyID: propertyID, InvestmentDate: today()}
	data.InvestorProperties = append(data.InvestorProperties, assignment)
	if !h.saveData(c, data) {
		return
	}
	c.JSON(http.StatusCreated, assignment)
}

func (h *Handler) UnsetInvestorProperty(c *gin.Context) {
	investorID := c.Param("id")
	propertyID := c.Param("propertyId")
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	data.InvestorProperties = filter(data.InvestorProperties, func(item dao.InvestorProperty) bool {
		return !(item.InvestorID == investorID && item.PropertyID == propertyID)
	})
	if !h.saveData(c, data) {
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ListUsers(c *gin.Context) {
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, data.Users)
}

func (h *Handler) GetUser(c *gin.Context) {
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	user, _ := findUser(data, c.Param("id"))
	if user == nil {
		notFound(c, "user not found")
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var payload dao.AdminUser
	if err := c.ShouldBindJSON(&payload); err != nil {
		badRequest(c, err.Error())
		return
	}
	if strings.TrimSpace(payload.Name) == "" || strings.TrimSpace(payload.Email) == "" {
		badRequest(c, "name and email are required")
		return
	}
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	payload.ID = newID("user")
	if payload.Role == "" {
		payload.Role = "STAFF"
	}
	if payload.CreatedAt == "" {
		payload.CreatedAt = nowISO()
	}
	data.Users = append(data.Users, payload)
	if !h.saveData(c, data) {
		return
	}
	c.JSON(http.StatusCreated, payload)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	var payload dao.AdminUser
	if err := c.ShouldBindJSON(&payload); err != nil {
		badRequest(c, err.Error())
		return
	}
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	user, _ := findUser(data, c.Param("id"))
	if user == nil {
		notFound(c, "user not found")
		return
	}
	payload.ID = user.ID
	if payload.CreatedAt == "" {
		payload.CreatedAt = user.CreatedAt
	}
	if payload.PasswordHash == "" {
		payload.PasswordHash = user.PasswordHash
	}
	*user = payload
	if !h.saveData(c, data) {
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	user, index := findUser(data, id)
	if index == -1 {
		notFound(c, "user not found")
		return
	}
	adminCount := 0
	for _, item := range data.Users {
		if item.Role == "ADMIN" {
			adminCount++
		}
	}
	if user.Role == "ADMIN" && adminCount <= 1 {
		badRequest(c, "at least one admin user is required")
		return
	}
	data.Users = append(data.Users[:index], data.Users[index+1:]...)
	for i := range data.Updates {
		if data.Updates[i].AuthorID == id {
			data.Updates[i].AuthorID = ""
		}
	}
	for i := range data.Documents {
		if data.Documents[i].UploadedBy == id {
			data.Documents[i].UploadedBy = ""
		}
	}
	if !h.saveData(c, data) {
		return
	}
	c.Status(http.StatusNoContent)
}

func filter[T any](items []T, keep func(T) bool) []T {
	next := items[:0]
	for _, item := range items {
		if keep(item) {
			next = append(next, item)
		}
	}
	return next
}

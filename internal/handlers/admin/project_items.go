package admin

import (
	"net/http"
	"strings"

	"github.com/conelli/admin-backend/internal/store/dao"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateUpdate(c *gin.Context) {
	var payload dao.Update
	if err := c.ShouldBindJSON(&payload); err != nil {
		badRequest(c, err.Error())
		return
	}
	if strings.TrimSpace(payload.PropertyID) == "" || strings.TrimSpace(payload.Title) == "" || strings.TrimSpace(payload.Body) == "" {
		badRequest(c, "propertyId, title, and body are required")
		return
	}
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	if property, _ := findProperty(data, payload.PropertyID); property == nil {
		notFound(c, "property not found")
		return
	}
	payload.ID = newID("upd")
	payload.PostedAt = nowISO()
	if payload.AuthorID == "" && len(data.Users) > 0 {
		payload.AuthorID = data.Users[0].ID
	}
	data.Updates = append([]dao.Update{payload}, data.Updates...)
	touchProperty(&data, payload.PropertyID)
	if !h.saveData(c, data) {
		return
	}
	c.JSON(http.StatusCreated, payload)
}

func (h *Handler) UpdateUpdate(c *gin.Context) {
	var payload dao.Update
	if err := c.ShouldBindJSON(&payload); err != nil {
		badRequest(c, err.Error())
		return
	}
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	for i := range data.Updates {
		if data.Updates[i].ID == c.Param("id") {
			payload.ID = data.Updates[i].ID
			if payload.PostedAt == "" {
				payload.PostedAt = data.Updates[i].PostedAt
			}
			data.Updates[i] = payload
			touchProperty(&data, payload.PropertyID)
			if !h.saveData(c, data) {
				return
			}
			c.JSON(http.StatusOK, data.Updates[i])
			return
		}
	}
	notFound(c, "update not found")
}

func (h *Handler) DeleteUpdate(c *gin.Context) {
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	for i := range data.Updates {
		if data.Updates[i].ID == c.Param("id") {
			propertyID := data.Updates[i].PropertyID
			data.Updates = append(data.Updates[:i], data.Updates[i+1:]...)
			touchProperty(&data, propertyID)
			if !h.saveData(c, data) {
				return
			}
			c.Status(http.StatusNoContent)
			return
		}
	}
	notFound(c, "update not found")
}

func (h *Handler) CreateMilestone(c *gin.Context) {
	var payload dao.Milestone
	if err := c.ShouldBindJSON(&payload); err != nil {
		badRequest(c, err.Error())
		return
	}
	if strings.TrimSpace(payload.PropertyID) == "" || strings.TrimSpace(payload.Title) == "" {
		badRequest(c, "propertyId and title are required")
		return
	}
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	payload.ID = newID("mile")
	if payload.Status == "" {
		payload.Status = "pending"
	}
	data.Milestones = append(data.Milestones, payload)
	touchProperty(&data, payload.PropertyID)
	if !h.saveData(c, data) {
		return
	}
	c.JSON(http.StatusCreated, payload)
}

func (h *Handler) UpdateMilestone(c *gin.Context) {
	var payload dao.Milestone
	if err := c.ShouldBindJSON(&payload); err != nil {
		badRequest(c, err.Error())
		return
	}
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	for i := range data.Milestones {
		if data.Milestones[i].ID == c.Param("id") {
			payload.ID = data.Milestones[i].ID
			data.Milestones[i] = payload
			touchProperty(&data, payload.PropertyID)
			if !h.saveData(c, data) {
				return
			}
			c.JSON(http.StatusOK, data.Milestones[i])
			return
		}
	}
	notFound(c, "milestone not found")
}

func (h *Handler) DeleteMilestone(c *gin.Context) {
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	for i := range data.Milestones {
		if data.Milestones[i].ID == c.Param("id") {
			propertyID := data.Milestones[i].PropertyID
			data.Milestones = append(data.Milestones[:i], data.Milestones[i+1:]...)
			touchProperty(&data, propertyID)
			if !h.saveData(c, data) {
				return
			}
			c.Status(http.StatusNoContent)
			return
		}
	}
	notFound(c, "milestone not found")
}

func (h *Handler) CreateDocument(c *gin.Context) {
	var payload dao.Document
	if err := c.ShouldBindJSON(&payload); err != nil {
		badRequest(c, err.Error())
		return
	}
	if strings.TrimSpace(payload.Title) == "" || strings.TrimSpace(payload.FileURL) == "" {
		badRequest(c, "title and fileUrl are required")
		return
	}
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	payload.ID = newID("doc")
	payload.UploadedAt = nowISO()
	if payload.UploadedBy == "" && len(data.Users) > 0 {
		payload.UploadedBy = data.Users[0].ID
	}
	data.Documents = append([]dao.Document{payload}, data.Documents...)
	if payload.PropertyID != nil {
		touchProperty(&data, *payload.PropertyID)
	}
	if !h.saveData(c, data) {
		return
	}
	c.JSON(http.StatusCreated, payload)
}

func (h *Handler) UpdateDocument(c *gin.Context) {
	var payload dao.Document
	if err := c.ShouldBindJSON(&payload); err != nil {
		badRequest(c, err.Error())
		return
	}
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	for i := range data.Documents {
		if data.Documents[i].ID == c.Param("id") {
			payload.ID = data.Documents[i].ID
			if payload.UploadedAt == "" {
				payload.UploadedAt = data.Documents[i].UploadedAt
			}
			if payload.UploadedBy == "" {
				payload.UploadedBy = data.Documents[i].UploadedBy
			}
			data.Documents[i] = payload
			if payload.PropertyID != nil {
				touchProperty(&data, *payload.PropertyID)
			}
			if !h.saveData(c, data) {
				return
			}
			c.JSON(http.StatusOK, data.Documents[i])
			return
		}
	}
	notFound(c, "document not found")
}

func (h *Handler) DeleteDocument(c *gin.Context) {
	data, ok := h.loadData(c)
	if !ok {
		return
	}
	for i := range data.Documents {
		if data.Documents[i].ID == c.Param("id") {
			propertyID := data.Documents[i].PropertyID
			data.Documents = append(data.Documents[:i], data.Documents[i+1:]...)
			if propertyID != nil {
				touchProperty(&data, *propertyID)
			}
			if !h.saveData(c, data) {
				return
			}
			c.Status(http.StatusNoContent)
			return
		}
	}
	notFound(c, "document not found")
}

// handlers/request_handlers.go
package handlers

import (
	"net/http"
	"strconv"
	"time"

	appdb "it-approval-backend/internal/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handlers struct {
	DB *gorm.DB
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	h := &Handlers{DB: db}

	r.GET("/health", h.Health)

	r.GET("/statuses", h.GetStatuses)

	r.GET("/requests", h.GetRequests)
	r.POST("/requests", h.CreateRequest)
	r.PATCH("/requests/:id/status", h.PatchRequestStatus)
}

func (h *Handlers) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GET /statuses
func (h *Handlers) GetStatuses(c *gin.Context) {
	var rows []appdb.Status
	if err := h.DB.Order("seq ASC").Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rows)
}

// GET /requests?limit=50&offset=0
func (h *Handlers) GetRequests(c *gin.Context) {
	limit := 50
	offset := 0

	if s := c.Query("limit"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 && v <= 200 {
			limit = v
		}
	}
	if s := c.Query("offset"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v >= 0 {
			offset = v
		}
	}

	var total int64
	if err := h.DB.Model(&appdb.Request{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var rows []appdb.Request
	if err := h.DB.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":  rows,
		"limit":  limit,
		"offset": offset,
		"total":  total,
	})
}

type CreateRequestBody struct {
	Title      string `json:"title" binding:"required"`
	StatusCode string `json:"status_code"` // optional; default "PENDING"
}

// POST /requests
func (h *Handlers) CreateRequest(c *gin.Context) {
	var body CreateRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	status := body.StatusCode
	if status == "" {
		status = "PENDING"
	}

	// optional: validate status exists
	var st appdb.Status
	if err := h.DB.First(&st, "code = ?", status).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status_code"})
		return
	}

	row := appdb.Request{
		Title:      body.Title,
		StatusCode: status,
		// created_at/updated_at ใช้ default/trigger จาก SQL ได้เลย
	}

	if err := h.DB.Create(&row).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// reload เพื่อให้ได้ created_at/updated_at ที่ DB set
	var out appdb.Request
	_ = h.DB.First(&out, "id = ?", row.ID).Error

	c.JSON(http.StatusOK, out)
}

type PatchStatusBody struct {
	StatusCode    *string `json:"status_code"`    // optional
	DecidedReason *string `json:"decided_reason"` // optional
	DecidedBy     *string `json:"decided_by"`     // optional
}

// PATCH /requests/:id/status
func (h *Handlers) PatchRequestStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var body PatchStatusBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	var row appdb.Request
	if err := h.DB.First(&row, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	updates := map[string]any{}

	if body.StatusCode != nil {
		// validate status exists
		var st appdb.Status
		if err := h.DB.First(&st, "code = ?", *body.StatusCode).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status_code"})
			return
		}
		updates["status_code"] = *body.StatusCode
		// ถือว่า “ตัดสินแล้ว” เมื่อมีการเปลี่ยนสถานะผ่าน endpoint นี้
		now := time.Now().UTC().Format(time.RFC3339)
		updates["decided_at"] = now
	}

	if body.DecidedReason != nil {
		updates["decided_reason"] = *body.DecidedReason
	}
	if body.DecidedBy != nil {
		updates["decided_by"] = *body.DecidedBy
	}

	if len(updates) == 0 {
		// ไม่ส่ง field มาเลย -> ไม่ทำอะไร แต่ตอบข้อมูลเดิม
		c.JSON(http.StatusOK, row)
		return
	}

	if err := h.DB.Model(&appdb.Request{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var out appdb.Request
	_ = h.DB.First(&out, "id = ?", id).Error
	c.JSON(http.StatusOK, out)
}

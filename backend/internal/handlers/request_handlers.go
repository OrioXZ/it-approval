// handlers/request_handlers.go
package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	appdb "it-approval-backend/internal/db"
	dbmodels "it-approval-backend/internal/db"

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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var body PatchStatusBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	// ⭐ transaction start
	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var row appdb.Request
	err = tx.First(&row, "id = ?", id).Error

	// ⭐ แยก not found vs db error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
		return
	}
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// ⭐ business rule: ถ้า status เดิมเป็น final → ห้ามแก้
	currentFinal, _ := isFinalStatus(tx, row.StatusCode)
	if currentFinal {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{
			"error": "cannot update request with final status",
		})
		return
	}

	updates := map[string]any{}

	// ⭐ validate status_code
	if body.StatusCode != nil {
		code := *body.StatusCode

		final, err := isFinalStatus(tx, code)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "unknown status_code"})
			return
		}

		// ถ้าเป็น final → ต้องมี reason + by
		if final {
			if body.DecidedReason == nil || body.DecidedBy == nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "final status requires decided_reason and decided_by",
				})
				return
			}

			updates["decided_at"] = time.Now().UTC().Format(time.RFC3339)
		}

		updates["status_code"] = code
	}

	if body.DecidedReason != nil {
		updates["decided_reason"] = *body.DecidedReason
	}
	if body.DecidedBy != nil {
		updates["decided_by"] = *body.DecidedBy
	}

	if len(updates) == 0 {
		tx.Rollback()
		c.JSON(http.StatusOK, row)
		return
	}

	if err := tx.Model(&row).Updates(updates).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	var out appdb.Request
	h.DB.First(&out, "id = ?", id)

	c.JSON(http.StatusOK, out)
}

func isFinalStatus(db *gorm.DB, code string) (bool, error) {
	var st dbmodels.Status
	if err := db.Where("code = ?", code).First(&st).Error; err != nil {
		return false, err
	}
	return st.IsFinal == "Y", nil
}

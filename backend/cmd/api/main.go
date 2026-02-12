package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	appdb "it-approval-backend/internal/db"

	"github.com/gin-gonic/gin"
)

func main() {
	_ = os.MkdirAll("db", 0755)

	gdb, err := appdb.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	if err := appdb.MigrateAndSeed(gdb, "migrations"); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/statuses", func(c *gin.Context) {
		var statuses []appdb.Status
		if err := gdb.Order("seq asc").Find(&statuses).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, statuses)
	})

	r.GET("/requests", func(c *gin.Context) {
		statusCode := c.Query("status_code")
		q := c.Query("q")

		limit := 50
		offset := 0
		if v := c.Query("limit"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 200 {
				limit = n
			}
		}
		if v := c.Query("offset"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n >= 0 {
				offset = n
			}
		}

		dbq := gdb.Model(&appdb.Request{})

		if statusCode != "" {
			dbq = dbq.Where("status_code = ?", statusCode)
		}
		if q != "" {
			dbq = dbq.Where("title LIKE ?", "%"+q+"%")
		}

		var total int64
		if err := dbq.Count(&total).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		var rows []appdb.Request
		if err := dbq.Order("created_at desc").Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"total":  total,
			"limit":  limit,
			"offset": offset,
			"data":   rows,
		})
	})

	r.POST("/requests", func(c *gin.Context) {
		var body struct {
			Title      string `json:"title"`
			StatusCode string `json:"status_code"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "invalid json"})
			return
		}
		if body.Title == "" {
			c.JSON(400, gin.H{"error": "title is required"})
			return
		}
		if body.StatusCode == "" {
			c.JSON(400, gin.H{"error": "status_code is required"})
			return
		}

		// validate status exists
		var st appdb.Status
		if err := gdb.First(&st, "code = ?", body.StatusCode).Error; err != nil {
			c.JSON(400, gin.H{"error": "invalid status_code"})
			return
		}

		req := appdb.Request{
			Title:      body.Title,
			StatusCode: body.StatusCode,
		}

		if err := gdb.Create(&req).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, req)
	})

	r.PATCH("/requests/:id/status", func(c *gin.Context) {
		id := c.Param("id")

		var body struct {
			StatusCode    string  `json:"status_code"`
			DecidedReason *string `json:"decided_reason"`
			DecidedBy     *string `json:"decided_by"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "invalid json"})
			return
		}
		if body.StatusCode == "" {
			c.JSON(400, gin.H{"error": "status_code is required"})
			return
		}

		// validate status exists
		var st appdb.Status
		if err := gdb.First(&st, "code = ?", body.StatusCode).Error; err != nil {
			c.JSON(400, gin.H{"error": "invalid status_code"})
			return
		}

		var req appdb.Request
		if err := gdb.First(&req, "id = ?", id).Error; err != nil {
			c.JSON(404, gin.H{"error": "request not found"})
			return
		}

		req.StatusCode = body.StatusCode
		req.DecidedReason = body.DecidedReason
		req.DecidedBy = body.DecidedBy

		now := time.Now().UTC().Format(time.RFC3339)
		req.DecidedAt = &now

		if err := gdb.Save(&req).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, req)
	})

	log.Println("listening on :8080")
	_ = r.Run(":8080")
}

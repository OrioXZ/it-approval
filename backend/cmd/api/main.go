package main

import (
	"log"
	"net/http"
	"os"

	appdb "it-approval-backend/internal/db"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	_ = os.MkdirAll("db", 0755)

	gdb, err := gorm.Open(sqlite.Open("db/app.db"), &gorm.Config{})
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

	log.Println("listening on :8080")
	_ = r.Run(":8080")
}

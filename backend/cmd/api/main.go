package main

import (
	"log"
	"os"
	"time"

	appdb "it-approval-backend/internal/db"
	"it-approval-backend/internal/handlers"

	"github.com/gin-contrib/cors"
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	handlers.RegisterRoutes(r, gdb)

	log.Println("listening on :8080")
	_ = r.Run(":8080")
}

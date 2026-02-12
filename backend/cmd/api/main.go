package main

import (
	"log"
	"os"

	appdb "it-approval-backend/internal/db"
	"it-approval-backend/internal/handlers"

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

	handlers.RegisterRoutes(r, gdb)

	log.Println("listening on :8080")
	_ = r.Run(":8080")
}

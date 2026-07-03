package main

import (
	"fmt"
	auth "healmata_backend/internal/auth"
	"healmata_backend/internal/db"
	"healmata_backend/internal/env"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 2. Read variables from environment
	cfg := env.Load()

	// 3. Initialize DB with the loaded configuration
	db, err := db.New(cfg.DBAddr, cfg.DBMaxOpenConn, cfg.DBMaxIdleConn, cfg.DBMaxIdleTime)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer db.Close()

	fmt.Println("Server starts on port 8080")

	r := gin.Default()
	auth.RegisterRoutes(r)
	r.Run(":" + cfg.Port)
}

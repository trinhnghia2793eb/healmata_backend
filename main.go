package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	auth "healmata_backend/internal/auth"
)

func main() {

	// config.ConnectDB()

	// config.DB.AutoMigrate(&models.User{})

	fmt.Println("Server starts on port 8080")

	r := gin.Default()

	auth.RegisterRoutes(r)
	r.Run(":8080")
}

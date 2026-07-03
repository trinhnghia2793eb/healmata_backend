package router

import (
    "github.com/gin-gonic/gin"

    "healmata_backend/internal/auth/handler"
)

func registerAuthRoutes(r *gin.Engine) {
    h := handler.NewAuthHandler()

    auth := r.Group("/auth")
    {
        auth.GET("/health", h.Health)
    }
}
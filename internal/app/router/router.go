package router

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
    registerAuthRoutes(r)

    // registerUserRoutes(r)
    // registerPaymentRoutes(r)
}
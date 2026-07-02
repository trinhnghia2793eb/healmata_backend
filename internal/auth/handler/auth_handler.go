package handler

import (
	// "healmata_backend/internal/auth/service"
	"github.com/gin-gonic/gin"
	// "errors"
)

type AuthHandler struct {
	// service service.AuthService
}

func NewAuthHandler() *AuthHandler {

	return &AuthHandler{
		// service: s,
	}
}

func (h *AuthHandler) Health(c *gin.Context) {
	// test error
	// var err error = errors.New("invalid credentials")

	// if err != nil {
	// 	c.JSON(401, gin.H{
	// 		"success": false,
	// 		"error": gin.H{
	// 			"code": "AUTH_ERROR_CODE",
	// 			"message": "ERROR_MESSAGE",
	// 		},
	// 	})
	// }

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"status": "ok",
		},
		"message": "AUTH_HEALTH_OK",
	})

}

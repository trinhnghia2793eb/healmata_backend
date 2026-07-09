package handler

import (
	"errors"
	"log"
	"healmata_backend/internal/auth/dto"
	authErrors "healmata_backend/internal/auth/errors"
	"healmata_backend/internal/auth/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {

	return &AuthHandler{
		service: s,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequestDTO

	// 1. Bind JSON and run basic DTO validation
	// Middleware will handle this with AUTH_REG_003, AUTH_REG_004
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "AUTH_REG_003",
				"message": "VALIDATION_FAILED: " + err.Error(),
			},
		})
		return
	}

	// 2. Call the service layer
	resp, err := h.service.Register(c, &req, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		var appErr *authErrors.AppError
		if errors.As(err, &appErr) {
			// If it's a known domain error, return its specific HTTP status and code
			c.JSON(appErr.HTTPStatus, gin.H{
				"success": false,
				"error": gin.H{
					"code":    appErr.ErrorCode,
					"message": appErr.Message,
				},
			})
			return
		}

		// Unexpected error, log it and return generic 500
		log.Printf("[Register] Unexpected error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "AUTH_REG_005",
				"message": "INTERNAL_ERROR",
			},
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    resp,
		"message": "REGISTER_SUCCESS",
	})
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

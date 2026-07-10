package handler

import (
	"errors"
	"healmata_backend/internal/auth/dto"
	authErrors "healmata_backend/internal/auth/errors"
	"healmata_backend/internal/auth/service"
	"log"
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
	// 1. Retrieve the validated DTO from the context (set by ValidateRegister middleware)
	reqVal, exists := c.Get("register_req")
	var req *dto.RegisterRequestDTO
	if exists {
		req = reqVal.(*dto.RegisterRequestDTO)
	} else {
		// Fallback in case middleware is missing
		var fallbackReq dto.RegisterRequestDTO
		if err := c.ShouldBindJSON(&fallbackReq); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"success": false,
				"error":   authErrors.NewAppError(http.StatusUnprocessableEntity, "AUTH_REG_003", "VALIDATION_FAILED: "+err.Error()),
			})
			return
		}
		req = &fallbackReq
	}

	// 2. Call the service layer
	resp, err := h.service.Register(c, req, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		var appErr *authErrors.AppError
		if errors.As(err, &appErr) {
			// If it's a known domain error, return its specific HTTP status and code
			c.JSON(appErr.HTTPStatus, gin.H{
				"success": false,
				"error":   appErr,
			})
			return
		}

		// Unexpected error, log it and return generic 500
		log.Printf("[Register] Unexpected error: %v", err)
		c.JSON(authErrors.ErrInternalError.HTTPStatus, gin.H{
			"success": false,
			"error":   authErrors.ErrInternalError,
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

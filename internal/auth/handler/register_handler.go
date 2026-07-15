package handler

import (
	"errors"
	"healmata_backend/internal/auth/dto"
	authError "healmata_backend/internal/auth/errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   authError.ErrInvalidJSON,
			})
			return
		}
		req = &fallbackReq
	}

	// 2. Call the service layer
	resp, err := h.service.Register(c, req, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		var appErr *authError.AppError
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   authError.ErrInternalError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
		"message": "REGISTER_SUCCESS",
	})
}

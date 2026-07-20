package handler

import (
	"errors"
	"healmata_backend/internal/auth/dto"
	authErrors "healmata_backend/internal/auth/errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// register
func (h *AuthHandler) Register(c *gin.Context) {
	// 1. Retrieve the validated DTO from the context (set by ValidateRegister middleware)
	reqVal, exists := c.Get("register_req")
	var req *dto.RegisterRequestDTO
	if exists {
		req = reqVal.(*dto.RegisterRequestDTO)
	} else {
		var fallbackReq dto.RegisterRequestDTO
		if err := c.ShouldBindJSON(&fallbackReq); err != nil {
			authErrors.ReturnAppError(c, validationErr.InvalidJson)
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
			authErrors.ReturnAppError(c, appErr)
			return
		}

		log.Printf("[Register] Unexpected error: %v", err)
		authErrors.ReturnAppError(c, validationErr.InternalError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
		"message": "REGISTER_SUCCESS",
	})
}

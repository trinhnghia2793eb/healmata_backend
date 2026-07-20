package handler

import (
	"errors"
	"healmata_backend/internal/auth/dto"
	authErrors "healmata_backend/internal/auth/errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login handles the POST /v1/auth/login API request
func (h *AuthHandler) Login(c *gin.Context) {
	// 1. Retrieve the validated DTO from the context (set by ValidateLogin middleware)
	reqVal, exists := c.Get("login_req")
	var req *dto.LoginRequestDTO
	if exists {
		req = reqVal.(*dto.LoginRequestDTO)
	} else {
		var fallbackReq dto.LoginRequestDTO
		if err := c.ShouldBindJSON(&fallbackReq); err != nil {
			authErrors.ReturnAppError(c, validationErr.InvalidJson)
			return
		}
		req = &fallbackReq
	}

	// 2. Call the service layer
	resp, err := h.service.Login(c.Request.Context(), req, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		var appErr *authErrors.AppError
		if errors.As(err, &appErr) {
			authErrors.ReturnAppError(c, appErr)
			return
		}

		log.Printf("[Login] Unexpected error: %v", err)
		authErrors.ReturnAppError(c, validationErr.InternalError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
		"message": "LOGIN_SUCCESS",
	})
}

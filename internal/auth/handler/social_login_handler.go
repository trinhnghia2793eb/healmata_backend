package handler

import (
	"errors"
	"healmata_backend/internal/auth/dto"
	authError "healmata_backend/internal/auth/errors"
	authErrors "healmata_backend/internal/auth/errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) SocialLogin(c *gin.Context) {
	reqVal, exists := c.Get("social_login_req")
	var req *dto.SocialLoginRequestDTO
	if exists {
		req = reqVal.(*dto.SocialLoginRequestDTO)
	} else {
		var fallbackReq dto.SocialLoginRequestDTO
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
	resp, err := h.service.SocialLogin(c.Request.Context(), req, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		var appErr *authError.AppError
		if errors.As(err, &appErr) {
			c.JSON(appErr.HTTPStatus, gin.H{
				"success": false,
				"error":   appErr,
			})
			return
		}

		// Unexpected error
		log.Printf("[SocialLogin] Unexpected error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   authErrors.ErrInternalError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
		"message": "SOCIAL_LOGIN_SUCCESS",
	})

}

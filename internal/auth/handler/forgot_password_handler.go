package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"healmata_backend/internal/auth/dto"
	authErrors "healmata_backend/internal/auth/errors"
)

// forgot password
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	reqVal, exists := c.Get("forgot_password_req")
	var req *dto.ForgotPasswordRequestDTO
	if exists {
		req = reqVal.(*dto.ForgotPasswordRequestDTO)
	} else {
		var fallbackReq dto.ForgotPasswordRequestDTO
		if err := c.ShouldBindJSON(&fallbackReq); err != nil {
			authErrors.ReturnAppError(c, validationErr.InvalidJson)
			return
		}
		req = &fallbackReq
	}

	resp, err := h.service.ForgotPassword(c.Request.Context(), req)
	if err != nil {
		var appErr *authErrors.AppError
		if errors.As(err, &appErr) {
			authErrors.ReturnAppError(c, appErr)
			return
		}

		log.Printf("[ForgotPassword] Unexpected error: %v", err)
		authErrors.ReturnAppError(c, validationErr.InternalError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
		"message": "OTP_SENT",
	})
}

// verify reset OTP
func (h *AuthHandler) VerifyResetOtp(c *gin.Context) {
	reqVal, exists := c.Get("verify_otp_req")
	var req *dto.VerifyResetOtpRequestDTO
	if exists {
		req = reqVal.(*dto.VerifyResetOtpRequestDTO)
	} else {
		var fallbackReq dto.VerifyResetOtpRequestDTO
		if err := c.ShouldBindJSON(&fallbackReq); err != nil {
			authErrors.ReturnAppError(c, validationErr.InvalidJson)
			return
		}
		req = &fallbackReq
	}

	resp, err := h.service.VerifyResetOtp(c.Request.Context(), req)
	if err != nil {
		var appErr *authErrors.AppError
		if errors.As(err, &appErr) {
			authErrors.ReturnAppError(c, appErr)
			return
		}

		log.Printf("[VerifyResetOtp] Unexpected error: %v", err)
		authErrors.ReturnAppError(c, validationErr.InternalError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
		"message": "OTP_VERIFIED",
	})
}

// reset user password
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	reqVal, exists := c.Get("reset_password_req")
	var req *dto.ResetPasswordRequestDTO

	if exists {
		req = reqVal.(*dto.ResetPasswordRequestDTO)
	} else {
		var fallbackReq dto.ResetPasswordRequestDTO
		if err := c.ShouldBindJSON(&fallbackReq); err != nil {
			authErrors.ReturnAppError(c, validationErr.InvalidJson)
			return
		}
		req = &fallbackReq
	}

	resp, err := h.service.ResetPassword(c.Request.Context(), req)
	if err != nil {
		var appErr *authErrors.AppError
		if errors.As(err, &appErr) {
			authErrors.ReturnAppError(c, appErr)
			return
		}

		log.Printf("[ResetPassword] Unexpected error: %v", err)
		authErrors.ReturnAppError(c, validationErr.InternalError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
		"message": "PASSWORD_RESET",
	})
}

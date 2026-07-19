package handler

import (
	"log"
	"net/http"
	"errors"
	
	"github.com/gin-gonic/gin"
	
	"healmata_backend/internal/auth/dto"
	authErrors "healmata_backend/internal/auth/errors"

)

// forgot password
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	// 1. Lấy DTO đã được validate từ context (thiết lập bởi middleware ValidateForgotPassword)
	reqVal, exists := c.Get("forgot_password_req")
	var req *dto.ForgotPasswordRequestDTO
	if exists {
		req = reqVal.(*dto.ForgotPasswordRequestDTO)
	} else {
		// Dự phòng (Fallback) trong trường hợp route quên gắn middleware
		var fallbackReq dto.ForgotPasswordRequestDTO
		if err := c.ShouldBindJSON(&fallbackReq); err != nil {
			c.JSON(authErrors.AUTH_FORGOT_003.HTTPStatus, gin.H{
				"success": false,
				"error": authErrors.ErrForgotInternal,
			})
			return
		}
		req = &fallbackReq
	}

	// 2. Gọi Service xử lý lỗi nghiệp vụ
	resp, err := h.service.ForgotPassword(c.Request.Context(), req)
	if err != nil {
		var appErr *authErrors.AppError
		// Nếu là lỗi nghiệp vụ (AppError) đã định nghĩa trong hệ thống
		if errors.As(err, &appErr) {
			c.JSON(appErr.HTTPStatus, gin.H{
				"success": false,
				"error": appErr,
			})
			return
		}

		// Lỗi hệ thống ngoài dự kiến (Unexpected error)
		log.Printf("[ForgotPassword] Unexpected error: %v", err)
		c.JSON(authErrors.ErrForgotInternal.HTTPStatus, gin.H{
			"success": false,
			"error": authErrors.ErrForgotInternal,
		})
		return
	}

	// 3. Trả về thành công với định dạng chuẩn
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
			c.JSON(authErrors.AUTH_OTP_004.HTTPStatus, gin.H{
				"success": false,
				"error": authErrors.ErrOtpInternal,
			})
			return
		}
		req = &fallbackReq
	}

	resp, err := h.service.VerifyResetOtp(c.Request.Context(), req)
	if err != nil {
		var appErr *authErrors.AppError
		if errors.As(err, &appErr) {
			c.JSON(appErr.HTTPStatus, gin.H{
				"success": false,
				"error":   appErr,
			})
			return
		}

		log.Printf("[VerifyResetOtp] Unexpected error: %v", err)
		c.JSON(authErrors.ErrOtpInternal.HTTPStatus, gin.H{
			"success": false,
			"error": authErrors.ErrOtpInternal,
		})
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
			c.JSON(authErrors.AUTH_RESET_004.HTTPStatus, gin.H{
				"success": false,
				"error":   authErrors.ErrResetPassInternal,
			})
			return
		}
		req = &fallbackReq
	}

	resp, err := h.service.ResetPassword(c.Request.Context(), req)
	if err != nil {
		var appErr *authErrors.AppError
		if errors.As(err, &appErr) {
			c.JSON(appErr.HTTPStatus, gin.H{
				"success": false,
				"error":   appErr,
			})
			return
		}

		log.Printf("[ResetPassword] Unexpected error: %v", err)
		c.JSON(authErrors.ErrResetPassInternal.HTTPStatus, gin.H{
			"success": false,
			"error":   authErrors.ErrResetPassInternal,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
		"message": "PASSWORD_RESET",
	})
}

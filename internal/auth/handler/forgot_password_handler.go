package handler

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"

	"healmata_backend/internal/auth/dto"
	"healmata_backend/pkg/response"
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
			response.Error(c, validationErr.InvalidJson)
			return
		}
		req = &fallbackReq
	}

	resp, err := h.service.ForgotPassword(c.Request.Context(), req)
	if err != nil {
		var appErr *response.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}

		log.Printf("[ForgotPassword] Unexpected error: %v", err)
		response.Error(c, validationErr.InternalError)
		return
	}

	response.Success(c, resp, "OTP_SENT")
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
			response.Error(c, validationErr.InvalidJson)
			return
		}
		req = &fallbackReq
	}

	resp, err := h.service.VerifyResetOtp(c.Request.Context(), req)
	if err != nil {
		var appErr *response.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}

		log.Printf("[VerifyResetOtp] Unexpected error: %v", err)
		response.Error(c, validationErr.InternalError)
		return
	}

	response.Success(c, resp, "OTP_VERIFIED")
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
			response.Error(c, validationErr.InvalidJson)
			return
		}
		req = &fallbackReq
	}

	resp, err := h.service.ResetPassword(c.Request.Context(), req)
	if err != nil {
		var appErr *response.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}

		log.Printf("[ResetPassword] Unexpected error: %v", err)
		response.Error(c, validationErr.InternalError)
		return
	}

	response.Success(c, resp, "PASSWORD_RESET")
}

package middleware

import (
	"strings"

	"healmata_backend/internal/auth/dto"
	authErrors "healmata_backend/internal/auth/errors"

	"github.com/gin-gonic/gin"
)

var validationErr = authErrors.Validation
var registerErr = authErrors.Register
var verifyOtpErr = authErrors.VerifyOtp
var resetPasswordErr = authErrors.ResetPassword

// register
func ValidateRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.RegisterRequestDTO

		if err := c.ShouldBindJSON(&req); err != nil {
			authErrors.ReturnValidationError(c, err)
			return
		}

		req.FullName = strings.TrimSpace(req.FullName)
		req.Identifier = strings.TrimSpace(req.Identifier)

		c.Set("register_req", &req)
		c.Next()
	}
}

// login
func ValidateLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.LoginRequestDTO

		if err := c.ShouldBindJSON(&req); err != nil {
			authErrors.ReturnValidationError(c, err)
			return
		}

		req.Identifier = strings.TrimSpace(req.Identifier)

		c.Set("login_req", &req)
		c.Next()
	}
}

// forgot password
func ValidateForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ForgotPasswordRequestDTO

		if err := c.ShouldBindJSON(&req); err != nil {
			authErrors.ReturnValidationError(c, err)
			return
		}

		req.Identifier = strings.TrimSpace(req.Identifier)

		c.Set("forgot_password_req", &req)
		c.Next()
	}
}

// verify reset OTP
func ValidateVerifyResetOtp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.VerifyResetOtpRequestDTO

		if err := c.ShouldBindJSON(&req); err != nil {
			authErrors.ReturnValidationError(c, err)
			return
		}

		req.ResetRequestId = strings.TrimSpace(req.ResetRequestId)
		req.Otp = strings.TrimSpace(req.Otp)

		c.Set("verify_otp_req", &req)
		c.Next()
	}
}

// reset password
func ValidateResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ResetPasswordRequestDTO

		if err := c.ShouldBindJSON(&req); err != nil {
			authErrors.ReturnValidationError(c, err)
			return
		}

		req.ResetToken = strings.TrimSpace(req.ResetToken)
		req.NewPassword = strings.TrimSpace(req.NewPassword)
		req.ConfirmPassword = strings.TrimSpace(req.ConfirmPassword)

		c.Set("reset_password_req", &req)
		c.Next()
	}
}

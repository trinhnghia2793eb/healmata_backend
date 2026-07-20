package middleware

import (
	"errors"
	"strings"

	"healmata_backend/internal/auth/dto"
	authErrors "healmata_backend/internal/auth/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
			var ve validator.ValidationErrors

			if errors.As(err, &ve) {
				for _, fe := range ve {
					switch fe.Field() {
					case "fullName":
						authErrors.ReturnAppError(c, registerErr.InvalidName)
						c.Abort()
						return
					case "password":
						authErrors.ReturnAppError(c, registerErr.InvalidPasswordReg)
						c.Abort()
						return
					case "confirmPassword":
						authErrors.ReturnAppError(c, registerErr.ConfirmPasswordMismatch)
						c.Abort()
						return
					case "identifier":
						if strings.Contains(req.Identifier, "@") {
							authErrors.ReturnAppError(c, validationErr.InvalidEmail)
						} else {
							authErrors.ReturnAppError(c, validationErr.InvalidPhone)
						}
						c.Abort()
						return
					}
				}
			}

			authErrors.ReturnAppError(c, validationErr.InvalidJson)
			c.Abort()
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
			var ve validator.ValidationErrors

			if errors.As(err, &ve) {
				for _, fe := range ve {
					switch fe.Field() {
					case "identifier":
						if strings.Contains(req.Identifier, "@") {
							authErrors.ReturnAppError(c, validationErr.InvalidEmail)
						} else {
							authErrors.ReturnAppError(c, validationErr.InvalidPhone)
						}
						c.Abort()
						return
					case "password":
						authErrors.ReturnAppError(c, validationErr.InvalidPassword)
						c.Abort()
						return
					}
				}
			}

			authErrors.ReturnAppError(c, validationErr.InvalidJson)
			c.Abort()
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

		// Bind JSON và kích hoạt bộ validator
		if err := c.ShouldBindJSON(&req); err != nil {
			var ve validator.ValidationErrors

			// Nếu lỗi do thư viện validator bắt được
			if errors.As(err, &ve) {
				for _, fe := range ve {
					switch fe.Field() {
					case "identifier":
						if strings.Contains(req.Identifier, "@") {
							authErrors.ReturnAppError(c, validationErr.InvalidEmail)
						} else {
							authErrors.ReturnAppError(c, validationErr.InvalidPhone)
						}
						c.Abort()
						return
					}
				}
			}

			authErrors.ReturnAppError(c, validationErr.InvalidJson)
			c.Abort()
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
			var ve validator.ValidationErrors

			if errors.As(err, &ve) {
				for _, fe := range ve {
					switch fe.Field() {
					case "Otp":
						authErrors.ReturnAppError(c, verifyOtpErr.InvalidOtp)
						c.Abort()
						return
					case "ResetRequestId":
						authErrors.ReturnAppError(c, verifyOtpErr.InternalError)
						c.Abort()
						return
					}
				}
			}

			authErrors.ReturnAppError(c, validationErr.InvalidJson)
			c.Abort()
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
			var ve validator.ValidationErrors

			if errors.As(err, &ve) {
				for _, fe := range ve {
					switch fe.Field() {
					case "newPassword":
						authErrors.ReturnAppError(c, resetPasswordErr.PasswordInvalid)
						c.Abort()
						return
					case "confirmPassword":
						authErrors.ReturnAppError(c, resetPasswordErr.PasswordMismatch)
						c.Abort()
						return
					case "resetToken":
						authErrors.ReturnAppError(c, resetPasswordErr.ResetTokenExpired)
						c.Abort()
						return
					}
				}
			}

			authErrors.ReturnAppError(c, validationErr.InvalidJson)
			c.Abort()
			return
		}

		req.ResetToken = strings.TrimSpace(req.ResetToken)
		req.NewPassword = strings.TrimSpace(req.NewPassword)
		req.ConfirmPassword = strings.TrimSpace(req.ConfirmPassword)

		c.Set("reset_password_req", &req)
		c.Next()
	}
}

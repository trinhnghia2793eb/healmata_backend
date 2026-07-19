package middleware

import (
	"errors"
	"strings"

	"healmata_backend/internal/auth/dto"
	customErrors "healmata_backend/internal/auth/errors" // Trỏ tới file chứa NewValidationError

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.RegisterRequestDTO

		// ShouldBindJSON bây giờ sẽ chạy MỌI THỨ: từ required, min=8 cho đến is_fullname, is_identifier
		if err := c.ShouldBindJSON(&req); err != nil {
			var ve validator.ValidationErrors

			if errors.As(err, &ve) {
				for _, fe := range ve {
					switch fe.Field() {
					case "fullName":
						c.JSON(customErrors.ErrInvalidName.HTTPStatus, gin.H{
							"success": false,
							"error":   customErrors.ErrInvalidName,
						})
						c.Abort()
						return
					case "password":
						c.JSON(customErrors.ErrInvalidPassword.HTTPStatus, gin.H{
							"success": false,
							"error":   customErrors.ErrInvalidPassword,
						})
						c.Abort()
						return
					case "confirmPassword":
						c.JSON(customErrors.ErrConfirmPasswordMismatch.HTTPStatus, gin.H{
							"success": false,
							"error":   customErrors.ErrConfirmPasswordMismatch,
						})
						c.Abort()
						return
					case "identifier":
						if strings.Contains(req.Identifier, "@") {
							c.JSON(customErrors.ErrInvalidEmail.HTTPStatus, gin.H{
								"success": false,
								"error":   customErrors.ErrInvalidEmail,
							})
						} else {
							c.JSON(customErrors.ErrInvalidPhone.HTTPStatus, gin.H{
								"success": false,
								"error":   customErrors.ErrInvalidPhone,
							})
						}
						c.Abort()
						return
					}
				}
			}

			// Lỗi JSON sai cú pháp
			c.JSON(customErrors.ErrInvalidJSON.HTTPStatus, gin.H{
				"success": false,
				"error":   customErrors.ErrInvalidJSON,
			})
			c.Abort()
			return
		}

		// Chỉ TrimSpace khi mọi thứ đã hợp lệ tuyệt đối
		req.FullName = strings.TrimSpace(req.FullName)
		req.Identifier = strings.TrimSpace(req.Identifier)

		c.Set("register_req", &req)
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
							c.JSON(customErrors.ErrInvalidEmail.HTTPStatus, gin.H{
								"success": false,
								"error":   customErrors.ErrInvalidEmail,
							})
						} else {
							c.JSON(customErrors.ErrInvalidPhone.HTTPStatus, gin.H{
								"success": false,
								"error":   customErrors.ErrInvalidPhone,
							})
						}
						c.Abort()
						return
					}
				}
			}

			// Bắt lỗi gửi lên JSON sai cấu trúc (fallback từ customErrors đã có sẵn)
			c.JSON(customErrors.ErrInvalidJSON.HTTPStatus, gin.H{
				"success": false,
				"error":   customErrors.ErrInvalidJSON,
			})
			c.Abort()
			return
		}

		req.Identifier = strings.TrimSpace(req.Identifier)

		// Lưu dữ liệu sạch vào context để Handler sử dụng
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
						c.JSON(customErrors.ErrInvalidOtp.HTTPStatus, gin.H{
							"success": false,
							"error":   customErrors.ErrInvalidOtp,
						})
						c.Abort()
						return
					case "ResetRequestId":
						c.JSON(customErrors.ErrOtpInternal.HTTPStatus, gin.H{
							"success": false,
							"error": customErrors.ErrOtpInternal,
						})
						c.Abort()
						return
					}
				}
			}

			c.JSON(customErrors.ErrInvalidJSON.HTTPStatus, gin.H{
				"success": false,
				"error":   customErrors.ErrInvalidJSON,
			})
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
						c.JSON(customErrors.ErrPasswordInvalid.HTTPStatus, gin.H{
							"success": false,
							"error":   customErrors.ErrPasswordInvalid,
						})
						c.Abort()
						return
					case "confirmPassword":
						c.JSON(customErrors.ErrPasswordMismatch.HTTPStatus, gin.H{
							"success": false,
							"error":   customErrors.ErrPasswordMismatch,
						})
						c.Abort()
						return
					case "resetToken":
						c.JSON(customErrors.ErrResetTokenExpired.HTTPStatus, gin.H{
							"success": false,
							"error":   customErrors.ErrResetTokenExpired,
						})
						c.Abort()
						return
					}
				}
			}

			c.JSON(customErrors.ErrInvalidJSON.HTTPStatus, gin.H{
				"success": false,
				"error":   customErrors.ErrInvalidJSON,
			})
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
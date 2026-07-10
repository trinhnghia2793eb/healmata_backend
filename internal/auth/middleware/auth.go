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

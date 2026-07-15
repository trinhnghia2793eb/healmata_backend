package middleware

import (
	"errors"
	"strings"

	"healmata_backend/internal/auth/dto"
	customErrors "healmata_backend/internal/auth/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.RegisterRequestDTO

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
						c.JSON(customErrors.ErrInvalidPasswordReg.HTTPStatus, gin.H{
							"success": false,
							"error":   customErrors.ErrInvalidPasswordReg,
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

			c.JSON(customErrors.ErrInvalidJSON.HTTPStatus, gin.H{
				"success": false,
				"error":   customErrors.ErrInvalidJSON,
			})
			c.Abort()
			return
		}

		req.FullName = strings.TrimSpace(req.FullName)
		req.Identifier = strings.TrimSpace(req.Identifier)

		c.Set("register_req", &req)
		c.Next()
	}
}

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
					case "password":
						c.JSON(customErrors.ErrInvalidPassword.HTTPStatus, gin.H{
							"success": false,
							"error":   customErrors.ErrInvalidPassword,
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

		req.Identifier = strings.TrimSpace(req.Identifier)

		c.Set("login_req", &req)
		c.Next()
	}
}

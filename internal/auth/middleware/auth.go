package middleware

import "github.com/gin-gonic/gin"

func Auth() gin.HandlerFunc {
	// Kiem tra request API co day du parameter
	// Kiem tra email, phone
	//
	return func(c *gin.Context) {
		c.Next()
	}
}

package handler

import "github.com/gin-gonic/gin"

type Handler struct{}

func (h *AuthHandler) Health(c *gin.Context) {
	// test error
	// var err error = errors.New("invalid credentials")

	// if err != nil {
	// 	c.JSON(401, gin.H{
	// 		"success": false,
	// 		"error": gin.H{
	// 			"code": "AUTH_ERROR_CODE",
	// 			"message": "ERROR_MESSAGE",
	// 		},
	// 	})
	// }

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"status": "ok",
		},
		"message": "AUTH_HEALTH_OK",
	})

}

package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	
	"github.com/gin-gonic/gin"
	"healmata_backend/internal/app/logger"
	"healmata_backend/pkg/response"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Ghi nhận lỗi panic kèm theo chi tiết Stack Trace xuống hệ thống log
				logger.Log.Error("recovery",
					slog.Any("error", err),
					slog.String("stack", string(debug.Stack())),
				)

				// Return 500
				c.AbortWithStatusJSON(
					http.StatusInternalServerError, 
					response.NewErrorResponse("INTERNAL_SERVER_ERROR", "Lỗi hệ thống"),
				)
			}
		}()
		c.Next()
	}
}
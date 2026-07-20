package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"healmata_backend/internal/app/logger"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		c.Next()

		// Sau khi Handler đã chạy xong
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		// log
		logger.Log.Info("API Request",
			slog.Int("status", statusCode),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("ip", clientIP),
			slog.Duration("latency", latency),
		)
	}
}
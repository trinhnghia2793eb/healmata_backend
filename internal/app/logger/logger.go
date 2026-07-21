package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Log là instance toàn cục để gọi ở khắp nơi trong ứng dụng
var Log *slog.Logger

// InitLogger cấu hình hệ thống log dựa trên môi trường chạy (dev hoặc prod)
func InitLogger(env string) {
	var handler slog.Handler

	// Chuẩn hóa chuỗi cấu hình môi trường
	env = strings.ToLower(strings.TrimSpace(env))

	if env == "production" || env == "prod" {
		// Ở Production: Log ra định dạng JSON để các hệ thống thu thập log (Loki/ELK) dễ parse
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo, // Chỉ log từ level INFO trở lên
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug, // Log tất cả bao gồm cả DEBUG
		})
	}

	// Khởi tạo instance
	Log = slog.New(handler)

	// Đặt làm mặc định cho package slog của hệ thống
	slog.SetDefault(Log)
}
package router

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"healmata_backend/internal/auth/handler"
	"healmata_backend/internal/auth/middleware"
	"healmata_backend/internal/auth/repository"
	"healmata_backend/internal/auth/service"
	"healmata_backend/internal/auth/token"
	"healmata_backend/internal/auth/validator"
	"healmata_backend/pkg/email"
)

func registerAuthRoutes(r *gin.Engine, db *pgxpool.Pool, emailSender email.EmailSender) {
	// Register custom validators
	validator.RegisterCustomValidators()
	// 1. Initialize JWT Manager (load secret from environment or fallback)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "heal-mata-secret-key-change-me-in-production"
	}

	// Access token: 1 hour, Refresh token: 30 days
	jwtManager := token.NewJWTManager(jwtSecret, 1*time.Hour, 30*24*time.Hour)

	// 2. Initialize repository
	repository := repository.NewAuthRepository(db)

	// 3. Initialize service
	authService := service.NewAuthService(repository, db, jwtManager, emailSender)

	// 4. Initialize handler
	h := handler.NewAuthHandler(authService)

	auth := r.Group("/auth")
	{
		auth.GET("/health", h.Health)
	}

	v1 := r.Group("/v1")
	{
		v1Auth := v1.Group("/auth")
		v1Auth.POST("/register", middleware.ValidateRegister(), h.Register)
		v1Auth.POST("/login", middleware.ValidateLogin(), h.Login)
		v1Auth.POST("/forgot-password", middleware.ValidateForgotPassword(), h.ForgotPassword)
		v1Auth.POST("/verify-reset-otp", middleware.ValidateVerifyResetOtp(), h.VerifyResetOtp)
		v1Auth.POST("/reset-password", middleware.ValidateResetPassword(), h.ResetPassword)
	}
}

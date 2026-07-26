package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"healmata_backend/internal/auth/handler"
	"healmata_backend/internal/auth/middleware"
	"healmata_backend/internal/auth/repository"
	"healmata_backend/internal/auth/service"
	"healmata_backend/internal/auth/validator"
	dbpkg "healmata_backend/pkg/db"
	"healmata_backend/pkg/email"
	"healmata_backend/pkg/jwt"
)

func registerAuthRoutes(
	r *gin.Engine,
	db *pgxpool.Pool,
	transactor *dbpkg.SQLTxManager,
	emailSender *email.Sender,
	jwtManager *jwt.JWTManager,
) {
	// Register custom validators
	validator.RegisterCustomValidators()

	// initialize repository
	repository := repository.NewAuthRepository(db)

	// initialize service
	authService := service.NewAuthService(repository, transactor, jwtManager, emailSender)

	// initialize handler
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

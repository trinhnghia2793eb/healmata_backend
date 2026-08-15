package testutils

import (
	"healmata_backend/internal/auth/handler"
	"healmata_backend/internal/auth/middleware"
	"healmata_backend/internal/auth/repository"
	"healmata_backend/internal/auth/service"
	"healmata_backend/internal/auth/validator"
	dbpkg "healmata_backend/pkg/db"
	"healmata_backend/pkg/jwt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupTestRouter(pool *pgxpool.Pool, emailSender *MockEmailSender) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 1. Chạy Custom Validators giống như app thật
	validator.RegisterCustomValidators()

	// 2. Khởi tạo các Dependencies còn thiếu
	transactor := dbpkg.NewSQLTxManager(pool)
	jwtManager := jwt.NewJWTManager("test-secret-key", 1*time.Hour, 24*time.Hour)

	// Khởi tạo Repository và Service
	repo := repository.NewAuthRepository(pool)
	svc := service.NewAuthService(repo, transactor, jwtManager, emailSender)
	h := handler.NewAuthHandler(svc)

	// 3. Map Route kèm Middleware chuẩn
	v1 := r.Group("/v1")
	{
		v1Auth := v1.Group("/auth")
		v1Auth.POST("/register", middleware.ValidateRegister(), h.Register)
		v1Auth.POST("/login", middleware.ValidateLogin(), h.Login)
		v1Auth.POST("/forgot-password", middleware.ValidateForgotPassword(), h.ForgotPassword)
		v1Auth.POST("/verify-reset-otp", middleware.ValidateVerifyResetOtp(), h.VerifyResetOtp)
		v1Auth.POST("/reset-password", middleware.ValidateResetPassword(), h.ResetPassword)
	}

	return r
}

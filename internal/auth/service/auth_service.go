package service

import (
	"context"
	"healmata_backend/internal/auth/dto"
	authErrors "healmata_backend/internal/auth/errors"
	"healmata_backend/internal/auth/repository"
	"healmata_backend/internal/auth/token"
	"healmata_backend/pkg/email"

	"github.com/jackc/pgx/v5/pgxpool"
)

var registerErr = authErrors.Register
var loginErr = authErrors.Login
var forgotPasswordErr = authErrors.ForgotPassword
var verifyOtpErr = authErrors.VerifyOtp
var resetPasswordErr = authErrors.ResetPassword

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequestDTO, clientIP, userAgent string) (*dto.RegisterResponseDTO, error)
	Login(ctx context.Context, req *dto.LoginRequestDTO, clientIP, userAgent string) (*dto.LoginResponseDTO, error)
	ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequestDTO) (*dto.ForgotPasswordResponseDTO, error)
	VerifyResetOtp(ctx context.Context, req *dto.VerifyResetOtpRequestDTO) (*dto.VerifyResetOtpResponseDTO, error)
	ResetPassword(ctx context.Context, req *dto.ResetPasswordRequestDTO) (*dto.ResetPasswordResponseDTO, error)
}

type authService struct {
	repo        repository.AuthRepository
	dbPool      *pgxpool.Pool
	jwtManager  *token.JWTManager
	emailSender email.EmailSender
}

func NewAuthService(
	repo repository.AuthRepository,
	dbPool *pgxpool.Pool,
	jwtManager *token.JWTManager,
	emailSender email.EmailSender,
) AuthService {
	return &authService{
		repo:        repo,
		dbPool:      dbPool,
		jwtManager:  jwtManager,
		emailSender: emailSender,
	}
}

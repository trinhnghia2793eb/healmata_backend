package service

import (
	"context"
	"healmata_backend/internal/auth/dto"
	authErrors "healmata_backend/internal/auth/errors"
	"healmata_backend/internal/auth/repository"
	dbpkg "healmata_backend/pkg/db"
	"healmata_backend/pkg/email"
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

type TokenProvider interface {
	GenerateAccessAndRefreshToken(userID string) (string, string, string, int64, error)
}

type authService struct {
	repo          repository.AuthRepository
	transactor    dbpkg.Transactor
	tokenProvider TokenProvider
	emailSender   email.EmailSender
}

func NewAuthService(
	repo repository.AuthRepository,
	transactor dbpkg.Transactor,
	tokenProvider TokenProvider,
	emailSender email.EmailSender,
) AuthService {
	return &authService{
		repo:          repo,
		transactor:    transactor,
		tokenProvider: tokenProvider,
		emailSender:   emailSender,
	}
}

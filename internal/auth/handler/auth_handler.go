package handler

import (
	"context"
	"healmata_backend/internal/auth/dto"
	authErrors "healmata_backend/internal/auth/errors"
)

var validationErr = authErrors.Validation

// =======================================================================
// required method from service for handler
type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequestDTO, clientIP, userAgent string) (*dto.RegisterResponseDTO, error)
	Login(ctx context.Context, req *dto.LoginRequestDTO, clientIP, userAgent string) (*dto.LoginResponseDTO, error)
	ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequestDTO) (*dto.ForgotPasswordResponseDTO, error)
	VerifyResetOtp(ctx context.Context, req *dto.VerifyResetOtpRequestDTO) (*dto.VerifyResetOtpResponseDTO, error)
	ResetPassword(ctx context.Context, req *dto.ResetPasswordRequestDTO) (*dto.ResetPasswordResponseDTO, error)
}
// =======================================================================

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(s AuthService) *AuthHandler {
	return &AuthHandler{
		service: s,
	}
}

package service

import (
	"context"
	authErrors "healmata_backend/internal/auth/errors"
	"healmata_backend/internal/auth/repository"
)

var registerErr = authErrors.Register
var loginErr = authErrors.Login
var forgotPasswordErr = authErrors.ForgotPassword
var verifyOtpErr = authErrors.VerifyOtp
var resetPasswordErr = authErrors.ResetPassword

// =======================================================================
// consumer defines the interfaces
type Transactor interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type EmailSender interface {
	SendOTP(to string, otpCode string) error
}

type TokenProvider interface {
	GenerateAccessAndRefreshToken(userID string) (string, string, string, int64, error)
}

// =======================================================================

type authService struct {
	repo          repository.AuthRepository
	transactor    Transactor
	tokenProvider TokenProvider
	emailSender   EmailSender
}

func NewAuthService(
	repo repository.AuthRepository,
	transactor Transactor,
	tokenProvider TokenProvider,
	emailSender EmailSender,
) *authService {
	return &authService{
		repo:          repo,
		transactor:    transactor,
		tokenProvider: tokenProvider,
		emailSender:   emailSender,
	}
}

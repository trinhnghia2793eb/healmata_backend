package service

import (
	"context"
	"healmata_backend/internal/auth/dto"
	"healmata_backend/internal/auth/repository"
	"healmata_backend/internal/auth/token"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequestDTO, clientIP, userAgent string) (*dto.RegisterResponseDTO, error)
	Login(ctx context.Context, req *dto.LoginRequestDTO, clientIP, userAgent string) (*dto.LoginResponseDTO, error)
	SocialLogin(ctx context.Context, req *dto.SocialLoginRequestDTO, clientIP, userAgent string) (*dto.SocialLoginResponseDTO, error)
}

type authService struct {
	repo       repository.AuthRepository
	dbPool     *pgxpool.Pool
	jwtManager *token.JWTManager
}

func NewAuthService(
	repo repository.AuthRepository,
	dbPool *pgxpool.Pool,
	jwtManager *token.JWTManager,
) AuthService {
	return &authService{
		repo:       repo,
		dbPool:     dbPool,
		jwtManager: jwtManager,
	}
}

func (s *authService) GenerateSessionAndTokens(ctx context.Context, tx pgx.Tx, userID string, clientIP, userAgent string) (*dto.TokenResponseDTO, error) {
	var response *dto.TokenResponseDTO
	// A. Create JWTs
	accessToken, rawRefreshToken, hashedRefreshToken, expiresIn, err := s.jwtManager.GenerateAccessAndRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	// B. Create Refresh Token record
	tokenPayload := &repository.CreateRefreshTokenPayload{
		UserID:    userID,
		TokenHash: hashedRefreshToken,
		DeviceID:  "default-device",
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}
	dbToken, err := s.repo.CreateRefreshToken(ctx, tx, tokenPayload)
	if err != nil {
		return nil, err
	}

	// C. Create Session record
	sessionPayload := &repository.CreateUserSessionPayload{
		UserID:         userID,
		RefreshTokenID: dbToken.ID,
		DeviceID:       "default-device",
		Platform:       "web",
		IPAddress:      clientIP,
		UserAgent:      userAgent,
	}
	if _, err := s.repo.CreateSession(ctx, tx, sessionPayload); err != nil {
		return nil, err
	}

	response = &dto.TokenResponseDTO{AccessToken: accessToken, RefreshToken: rawRefreshToken, ExpiresIn: expiresIn}
	return response, nil
}

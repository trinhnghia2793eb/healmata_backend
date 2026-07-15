package service

import (
	"context"
	"healmata_backend/internal/auth/dto"
	"healmata_backend/internal/auth/repository"
	"healmata_backend/internal/auth/token"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequestDTO, clientIP, userAgent string) (*dto.RegisterResponseDTO, error)
	Login(ctx context.Context, req *dto.LoginRequestDTO, clientIP, userAgent string) (*dto.LoginResponseDTO, error)
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

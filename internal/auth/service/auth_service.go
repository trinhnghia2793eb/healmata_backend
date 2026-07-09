package service

import (
	"context"
	"healmata_backend/internal/auth/dto"
	authError "healmata_backend/internal/auth/errors"
	"healmata_backend/internal/auth/model"
	"healmata_backend/internal/auth/repository"
	"healmata_backend/internal/auth/token"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequestDTO, clientIP, userAgent string) (*dto.RegisterResponseDTO, error)
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

func (s *authService) Register(ctx context.Context, req *dto.RegisterRequestDTO, clientIP, userAgent string) (*dto.RegisterResponseDTO, error) {
	// 1. Identify input type (Email vs Phone)
	var email, phone string
	isEmail := strings.Contains(req.Identifier, "@")
	if isEmail {
		email = req.Identifier
	} else {
		phone = req.Identifier
	}

	// 2. Check if user already exists
	existingUser, err := s.repo.GetUserByIdentifier(ctx, req.Identifier)
	if err == nil && existingUser != nil {
		if isEmail && existingUser.Email != "" {
			return nil, authError.NewAppError(http.StatusConflict, "AUTH_REG_001", authError.ErrEmailExists.Error())
		} else if !isEmail && existingUser.Phone != "" {
			return nil, authError.NewAppError(http.StatusConflict, "AUTH_REG_002", authError.ErrPhoneExists.Error())
		}
	}

	// 3. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, authError.NewAppError(http.StatusInternalServerError, "AUTH_REG_003", authError.ErrInternalError.Error())
	}

	// 4. Build user model
	user := &model.User{
		FullName:            req.FullName,
		Email:               email,
		Phone:               phone,
		PasswordHash:        string(hashedPassword),
		Status:              "active",
		FirstSetupCompleted: false,
	}

	// 5. Run writes inside a transaction

	return nil, nil
}

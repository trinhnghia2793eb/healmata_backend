package service

import (
	"context"
	"errors"
	"healmata_backend/internal/auth/dto"
	authErrors "healmata_backend/internal/auth/errors"
	"healmata_backend/internal/auth/repository"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *authService) Login(ctx context.Context, req *dto.LoginRequestDTO, clientIP, userAgent string) (*dto.LoginResponseDTO, error) {
	// 1. Fetch user by identifier
	user, err := s.repo.GetUserByIdentifier(ctx, req.Identifier)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, loginErr.UserNotFound
		}
		return nil, loginErr.InternalError
	}
	// 2. Verify account is not disable
	if user.Status != "active" {
		return nil, loginErr.UserDisabled
	}
	// 3. Verify bcrypt password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, loginErr.InvalidCredential
	}

	var response *dto.LoginResponseDTO
	// 4. Generate JWT tokens and persist sessions inside a database transaction
	err = s.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		// A. Create JWTs
		accessToken, rawRefreshToken, hashedRefreshToken, expiresIn, err := s.tokenProvider.GenerateAccessAndRefreshToken(user.ID)
		if err != nil {
			return err
		}

		// C. Create Refresh Token record
		tokenPayload := &repository.CreateRefreshTokenPayload{
			UserID:    user.ID,
			TokenHash: hashedRefreshToken,
			DeviceID:  "default-device",
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		}
		dbToken, err := s.repo.CreateRefreshToken(txCtx, tokenPayload)
		if err != nil {
			return err
		}

		// D. Create Session record
		sessionPayload := &repository.CreateUserSessionPayload{
			UserID:         user.ID,
			RefreshTokenID: dbToken.ID,
			DeviceID:       "default-device",
			Platform:       "web",
			IPAddress:      clientIP,
			UserAgent:      userAgent,
		}
		if _, err := s.repo.CreateSession(txCtx, sessionPayload); err != nil {
			return err
		}

		response = &dto.LoginResponseDTO{AccessToken: accessToken, RefreshToken: rawRefreshToken, ExpiresIn: expiresIn}

		return nil
	})
	if err != nil {
		log.Printf("[Login Service] underlying error: %v", err)
		var appErr *authErrors.AppError
		if errors.As(err, &appErr) {
			return nil, appErr
		}
		return nil, loginErr.InternalError
	}
	return response, nil
}

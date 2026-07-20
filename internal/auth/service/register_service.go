package service

import (
	"context"
	"errors"
	"healmata_backend/internal/auth/dto"
	authErrors "healmata_backend/internal/auth/errors"
	"healmata_backend/internal/auth/repository"
	"healmata_backend/pkg/db"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

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
			return nil, registerErr.EmailExists
		} else if !isEmail && existingUser.Phone != "" {
			return nil, registerErr.PhoneExists
		}
	}

	// 3. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, registerErr.InternalError
	}

	// 4. Build user model
	userPayload := &repository.CreateUserPayload{
		FullName:     req.FullName,
		Email:        nilOrStringPtr(email),
		Phone:        nilOrStringPtr(phone),
		PasswordHash: string(hashedPassword),
	}

	var response *dto.RegisterResponseDTO
	// 5. Run writes inside a transaction
	err = db.WithTransaction(ctx, s.dbPool, func(tx pgx.Tx) error {
		// A. Create User record
		user, err := s.repo.CreateUser(ctx, tx, userPayload)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				switch pgErr.ConstraintName {
				case "users_email_unique_idx":
					return registerErr.EmailExists
				case "users_phone_key":
					return registerErr.PhoneExists
				}
			}
			return err
		}

		// B. Generate JWTs
		accessToken, rawRefreshToken, hashedRefreshToken, expiresIn, err := s.jwtManager.GenerateAccessAndRefreshToken(user.ID)
		if err != nil {
			return err
		}

		// C. Hash refresh token before saving (using SHA-256 because bcrypt has a 72-byte limit and refresh token is a JWT)
		// hasher := sha256.New()
		// hasher.Write([]byte(rawRefreshToken))
		// tokenHash := hex.EncodeToString(hasher.Sum(nil))

		tokenPayload := &repository.CreateRefreshTokenPayload{
			UserID:    user.ID,
			TokenHash: hashedRefreshToken,
			DeviceID:  "default-device",
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		}
		dbToken, err := s.repo.CreateRefreshToken(ctx, tx, tokenPayload)
		if err != nil {
			return err
		}

		// E. Store Session
		sessionPayload := &repository.CreateUserSessionPayload{
			UserID:         user.ID,
			RefreshTokenID: dbToken.ID,
			DeviceID:       "default-device",
			Platform:       "web",
			IPAddress:      clientIP,
			UserAgent:      userAgent,
		}

		if _, err := s.repo.CreateSession(ctx, tx, sessionPayload); err != nil {
			return err
		}

		response = &dto.RegisterResponseDTO{AccessToken: accessToken, RefreshToken: rawRefreshToken, ExpiresIn: expiresIn}

		return nil
	})

	if err != nil {
		log.Printf("[Register Service] underlying error: %v", err)
		var appErr *authErrors.AppError
		if errors.As(err, &appErr) {
			return nil, appErr
		}
		return nil, registerErr.InternalError
	}

	return response, nil
}

func nilOrStringPtr(val string) *string {
	if val == "" {
		return nil
	}
	return &val
}

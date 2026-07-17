package service

import (
	"context"
	"errors"
	"healmata_backend/internal/auth/dto"
	authError "healmata_backend/internal/auth/errors"
	"healmata_backend/internal/auth/model"
	"healmata_backend/internal/auth/providers"
	"healmata_backend/internal/auth/repository"
	"healmata_backend/pkg/db"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *authService) SocialLogin(ctx context.Context, req *dto.SocialLoginRequestDTO, clientIP, userAgent string) (*dto.SocialLoginResponseDTO, error) {
	// Phase 1: Verification
	providerUser, err := s.verifySocialToken(ctx, req)
	if err != nil {
		return nil, err
	}

	// Phase 2: State Resolution
	userID, isNewUser, checkSocialAccount, err := s.resolveSocialUserState(ctx, req.Provider, providerUser)
	if err != nil {
		return nil, err
	}

	// Phase 3: DB Execution (Transaction)
	var response *dto.SocialLoginResponseDTO
	err = db.WithTransaction(ctx, s.dbPool, func(tx pgx.Tx) error {
		// A. Register user if they don't exist
		if isNewUser {
			userPayload := &repository.CreateUserPayload{
				FullName:     providerUser.FullName,
				Email:        nilOrStringPtr(providerUser.Email),
				PasswordHash: "",
			}
			user, err := s.repo.CreateUser(ctx, tx, userPayload)
			if err != nil {
				return err
			}
			userID = user.ID
		}
		// B. Link social account if new
		if checkSocialAccount == nil {
			socialAccountPayload := &repository.CreateSocialAccountPayload{
				UserID:         userID,
				Provider:       req.Provider,
				ProviderUserID: providerUser.ProviderUserID,
				ProviderEmail:  providerUser.Email,
			}
			if _, err := s.repo.CreateSocialAccount(ctx, tx, socialAccountPayload); err != nil {
				return err
			}
		}
		// C. Issue tokens
		tokenResponse, err := s.GenerateSessionAndTokens(ctx, tx, userID, clientIP, userAgent)
		if err != nil {
			return err
		}
		response = &dto.SocialLoginResponseDTO{
			TokenResponseDTO: *tokenResponse,
			IsNewUser:        isNewUser,
			LinkedAccount:    true, // Unconditionally true on successful transaction mapping
		}
		return nil
	})

	// Phase 4: Production-Safe Error Mapping
	if err != nil {
		log.Printf("[SocialLogin Service] transaction failed: %v", err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "users_email_unique_idx" {
				return nil, authError.AUTH_SOCIAL_003 // Conflict / Already linked
			}
		}
		return nil, authError.AUTH_SOCIAL_004
	}
	return response, nil
}

func (s *authService) verifySocialToken(ctx context.Context, req *dto.SocialLoginRequestDTO) (*providers.ProviderUser, error) {
	var providerUser *providers.ProviderUser
	var err error

	switch req.Provider {
	case "google":
		clientIDs := []string{
			os.Getenv("GOOGLE_CLIENT_ID_WEB"),
			os.Getenv("GOOGLE_CLIENT_ID_IOS"),
			os.Getenv("GOOGLE_CLIENT_ID_ANDROID"),
		}
		providerUser, err = providers.VerifyGoogleToken(ctx, req.ProviderToken, clientIDs)
	case "apple":
		clientIDs := []string{
			os.Getenv("APPLE_CLIENT_ID_WEB"),
			os.Getenv("APPLE_CLIENT_ID_IOS"),
			os.Getenv("APPLE_CLIENT_ID_ANDROID"),
		}
		providerUser, err = providers.VerifyAppleToken(ctx, req.ProviderToken, clientIDs)
	default:
		return nil, authError.AUTH_SOCIAL_001
	}
	if err != nil {
		return nil, authError.AUTH_SOCIAL_002
	}
	return providerUser, nil
}

func (s *authService) resolveSocialUserState(ctx context.Context, provider string, providerUser *providers.ProviderUser) (string, bool, *model.SocialAccounts, error) {
	// 1. Check if social account is already mapped
	checkSocialAccount, err := s.repo.GetSocialAccount(ctx, provider, providerUser.ProviderUserID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", false, nil, authError.AUTH_SOCIAL_004 // INTERNAL_ERROR
	}
	if checkSocialAccount != nil {
		return checkSocialAccount.UserID, false, checkSocialAccount, nil // Exists under social account
	}
	// 2. Map by email if social account doesn't exist
	user, err := s.repo.GetUserByIdentifier(ctx, providerUser.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", false, nil, authError.AUTH_SOCIAL_004 // INTERNAL_ERROR
	}
	if user != nil {
		return user.ID, false, nil, nil // Exists under email, needs mapping
	}
	return "", true, nil, nil // Entirely new user
}

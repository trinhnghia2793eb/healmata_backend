package repository

import (
	"context"
	"healmata_backend/internal/auth/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	GetUserByIdentifier(ctx context.Context, identifier string) (*model.User, error)
	CreateUser(ctx context.Context, tx pgx.Tx, user *CreateUserPayload) (*model.User, error)
	CreateRefreshToken(ctx context.Context, tx pgx.Tx, refreshToken *CreateRefreshTokenPayload) (*model.RefreshTokens, error)
	CreateSession(ctx context.Context, tx pgx.Tx, session *CreateUserSessionPayload) (*model.UserSessions, error)
	GetSocialAccount(ctx context.Context, provider string, providerUserID string) (*model.SocialAccounts, error)
	CreateSocialAccount(ctx context.Context, tx pgx.Tx, socialAcc *CreateSocialAccountPayload) (*model.SocialAccounts, error)
}

type authRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepository {

	return &authRepository{db: db}
}

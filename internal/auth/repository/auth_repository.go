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
}

type authRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepository {

	return &authRepository{db: db}
}

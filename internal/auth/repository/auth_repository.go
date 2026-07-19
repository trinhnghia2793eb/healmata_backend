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

	CreateOtpRequest(ctx context.Context, tx pgx.Tx, payload *CreateOtpRequestPayload) (*model.OtpRequest, error)
	GetLatestOtpRequest(ctx context.Context, identifier string, purpose string) (*model.OtpRequest, error)

	GetOtpRequestByID(ctx context.Context, tx pgx.Tx, id string) (*model.OtpRequest, error)
	UpdateOtpRequest(ctx context.Context, tx pgx.Tx, otpReq *model.OtpRequest) error

	GetOtpRequestByTokenHash(ctx context.Context, tx pgx.Tx, tokenHash string) (*model.OtpRequest, error)
	UpdateUserPassword(ctx context.Context, tx pgx.Tx, identifier string, passwordHash string) error
	InvalidateResetToken(ctx context.Context, tx pgx.Tx, id string) error
}

type authRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepository {

	return &authRepository{db: db}
}

package repository

import (
	"context"
	"healmata_backend/internal/auth/model"
	dbpkg "healmata_backend/pkg/db"
)

type AuthRepository interface {
	GetUserByIdentifier(ctx context.Context, identifier string) (*model.User, error)
	CreateUser(ctx context.Context, user *CreateUserPayload) (*model.User, error)
	CreateRefreshToken(ctx context.Context, refreshToken *CreateRefreshTokenPayload) (*model.RefreshTokens, error)
	CreateSession(ctx context.Context, session *CreateUserSessionPayload) (*model.UserSessions, error)

	CreateOtpRequest(ctx context.Context, payload *CreateOtpRequestPayload) (*model.OtpRequest, error)
	GetLatestOtpRequest(ctx context.Context, identifier string, purpose string) (*model.OtpRequest, error)

	GetOtpRequestByID(ctx context.Context, id string) (*model.OtpRequest, error)
	UpdateOtpRequest(ctx context.Context, otpReq *model.OtpRequest) error

	GetOtpRequestByTokenHash(ctx context.Context, tokenHash string) (*model.OtpRequest, error)
	UpdateUserPassword(ctx context.Context, identifier string, passwordHash string) error
	InvalidateResetToken(ctx context.Context, id string) error
}

type authRepository struct {
	db dbpkg.DBEngine
}

func NewAuthRepository(db dbpkg.DBEngine) AuthRepository {
	return &authRepository{
		db: db,
	}
}

// if it has tx in ctx ? use tx transaction : use r.db
func (r *authRepository) getDB(ctx context.Context) dbpkg.DBEngine {
	if tx := dbpkg.GetTxFromContext(ctx); tx != nil {
		return tx
	}
	return r.db
}

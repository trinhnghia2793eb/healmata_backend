package repository

import (
	"context"
	"healmata_backend/internal/auth/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	GetUserByIdentifier(ctx context.Context, identifier string) (*model.User, error)
	CreateUser(ctx context.Context, tx pgx.Tx, user *model.User) error
	CreateRefreshToken(ctx context.Context, tx pgx.Tx, refreshToken *model.RefreshTokens) error
	CreateSession(ctx context.Context, tx pgx.Tx, session *model.UserSessions) error
}

type authRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepository {

	return &authRepository{db: db}
}

func (r *authRepository) GetUserByIdentifier(ctx context.Context, identifier string) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE LOWER(email) = LOWER($1) OR phone = $1`
	err := r.db.QueryRow(ctx, query, identifier).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.Status,
		&user.FirstSetupCompleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) CreateUser(ctx context.Context, tx pgx.Tx, user *model.User) error {
	query := `
		INSERT INTO users (full_name, email, phone, password_hash, status, first_setup_completed)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	err := tx.QueryRow(ctx, query,
		user.FullName,
		nilOrValue(user.Email),
		nilOrValue(user.Phone),
		user.PasswordHash,
		user.Status,
		user.FirstSetupCompleted,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

func (r *authRepository) CreateRefreshToken(ctx context.Context, tx pgx.Tx, refreshToken *model.RefreshTokens) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, device_id, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, expires_at, created_at
	`
	err := tx.QueryRow(ctx, query,
		refreshToken.UserID,
		refreshToken.TokenHash,
		refreshToken.DeviceID,
		refreshToken.ExpiresAt,
	).Scan(&refreshToken.ID, &refreshToken.ExpiresAt, &refreshToken.CreatedAt)

	return err
}

func (r *authRepository) CreateSession(ctx context.Context, tx pgx.Tx, session *model.UserSessions) error {
	query := `
		INSERT INTO user_sessions (user_id, refresh_token_id, device_id, platform, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	err := tx.QueryRow(ctx, query,
		session.UserID,
		session.RefreshTokenID,
		session.DeviceID,
		session.Platform,
		session.IPAdress,
		session.UserAgent,
	).Scan(&session.ID, &session.CreatedAt)

	return err
}

func nilOrValue(val string) *string {
	if val == "" {
		return nil
	}
	return &val
}

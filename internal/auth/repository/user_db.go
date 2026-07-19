package repository

import (
	"context"
	"healmata_backend/internal/auth/model"

	"github.com/jackc/pgx/v5"
)

func (r *authRepository) GetUserByIdentifier(ctx context.Context, identifier string) (*model.User, error) {
	var user model.User
	var emailPtr, phonePtr, passwordHashPtr *string
	query := `SELECT id, full_name, email, phone, password_hash, status, first_setup_completed, created_at, updated_at 
				FROM users WHERE LOWER(email) = LOWER($1) OR phone = $1`
	err := r.db.QueryRow(ctx, query, identifier).Scan(
		&user.ID,
		&user.FullName,
		&emailPtr,
		&phonePtr,
		&passwordHashPtr,
		&user.Status,
		&user.FirstSetupCompleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if emailPtr != nil {
		user.Email = *emailPtr
	}
	if phonePtr != nil {
		user.Phone = *phonePtr
	}
	if passwordHashPtr != nil {
		user.PasswordHash = *passwordHashPtr
	}
	return &user, nil
}

func (r *authRepository) CreateUser(ctx context.Context, tx pgx.Tx, payload *CreateUserPayload) (*model.User, error) {
	query := `
		INSERT INTO users (full_name, email, phone, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, status, first_setup_completed, created_at, updated_at
	`
	var user model.User
	err := tx.QueryRow(ctx, query,
		payload.FullName,
		payload.Email,
		payload.Phone,
		payload.PasswordHash,
	).Scan(
		&user.ID,
		&user.Status,
		&user.FirstSetupCompleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	user.FullName = payload.FullName
	if payload.Email != nil {
		user.Email = *payload.Email
	}
	if payload.Phone != nil {
		user.Phone = *payload.Phone
	}
	user.PasswordHash = payload.PasswordHash
	return &user, nil
}

func (r *authRepository) CreateRefreshToken(ctx context.Context, tx pgx.Tx, payload *CreateRefreshTokenPayload) (*model.RefreshTokens, error) {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, device_id, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	var token model.RefreshTokens
	err := tx.QueryRow(ctx, query,
		payload.UserID,
		payload.TokenHash,
		payload.DeviceID,
		payload.ExpiresAt,
	).Scan(&token.ID, &token.CreatedAt)

	if err != nil {
		return nil, err
	}

	token.UserID = payload.UserID
	token.TokenHash = payload.TokenHash
	token.DeviceID = payload.DeviceID
	token.ExpiresAt = payload.ExpiresAt
	return &token, nil
}

func (r *authRepository) CreateSession(ctx context.Context, tx pgx.Tx, payload *CreateUserSessionPayload) (*model.UserSessions, error) {
	query := `
		INSERT INTO user_sessions (user_id, refresh_token_id, device_id, platform, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	var session model.UserSessions
	err := tx.QueryRow(ctx, query,
		payload.UserID,
		payload.RefreshTokenID,
		payload.DeviceID,
		payload.Platform,
		payload.IPAddress,
		payload.UserAgent,
	).Scan(&session.ID, &session.CreatedAt)

	if err != nil {
		return nil, err
	}

	session.UserID = payload.UserID
	session.RefreshTokenID = payload.RefreshTokenID
	session.DeviceID = payload.DeviceID
	session.Platform = payload.Platform
	session.IPAdress = payload.IPAddress
	session.UserAgent = payload.UserAgent

	return &session, nil
}

func (r *authRepository) UpdateUserPassword(ctx context.Context, tx pgx.Tx, identifier string, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
		WHERE LOWER(email) = LOWER($2) OR phone = $2
	`
	_, err := tx.Exec(ctx, query, passwordHash, identifier)
	return err
}
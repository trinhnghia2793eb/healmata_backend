package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"healmata_backend/internal/auth/model"
)

func (r *authRepository) CreateOtpRequest(ctx context.Context, tx pgx.Tx, payload *CreateOtpRequestPayload) (*model.OtpRequest, error) {
	query := `
		INSERT INTO otp_requests (identifier, otp_hash, purpose, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, attempts, created_at
	`
	var otpReq model.OtpRequest
	// sử dụng tx truyền từ tầng Service xuống
	err := tx.QueryRow(ctx, query,
		payload.Identifier, payload.OtpHash, payload.Purpose, payload.ExpiresAt,
	).Scan(
		&otpReq.ID, &otpReq.Attempts, &otpReq.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	otpReq.Identifier = payload.Identifier
	otpReq.OtpHash = payload.OtpHash
	otpReq.Purpose = payload.Purpose
	otpReq.ExpiresAt = payload.ExpiresAt

	return &otpReq, nil
}

func (r *authRepository) GetLatestOtpRequest(ctx context.Context, identifier string, purpose string) (*model.OtpRequest, error) {
	query := `
		SELECT id, identifier, otp_hash, purpose, attempts, expires_at, verified_at, reset_token_hash, token_expires_at, created_at
		FROM otp_requests
		WHERE identifier = $1 AND purpose = $2
		ORDER BY created_at DESC
		LIMIT 1
	`
	var otpReq model.OtpRequest
	// Hàm đọc sử dụng r.db
	err := r.db.QueryRow(ctx, query, identifier, purpose).Scan(
		&otpReq.ID,
		&otpReq.Identifier,
		&otpReq.OtpHash,
		&otpReq.Purpose,
		&otpReq.Attempts,
		&otpReq.ExpiresAt,
		&otpReq.VerifiedAt,
		&otpReq.ResetTokenHash,
		&otpReq.TokenExpiresAt,
		&otpReq.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Không có lỗi, chỉ là chưa từng yêu cầu OTP
		}
		return nil, err
	}
	return &otpReq, nil
}

func (r *authRepository) GetOtpRequestByID(ctx context.Context, tx pgx.Tx, id string) (*model.OtpRequest, error) {
	// lấy thông tin OTP và khóa bản ghi (FOR UPDATE) --> tránh Race Condition
	query := `
		SELECT id, identifier, otp_hash, purpose, attempts, expires_at, verified_at, reset_token_hash, token_expires_at, created_at
		FROM otp_requests
		WHERE id = $1
		FOR UPDATE
	`
	var otpReq model.OtpRequest
	err := tx.QueryRow(ctx, query, id).Scan(
		&otpReq.ID,
		&otpReq.Identifier,
		&otpReq.OtpHash,
		&otpReq.Purpose,
		&otpReq.Attempts,
		&otpReq.ExpiresAt,
		&otpReq.VerifiedAt,
		&otpReq.ResetTokenHash,
		&otpReq.TokenExpiresAt,
		&otpReq.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &otpReq, nil
}

func (r *authRepository) UpdateOtpRequest(ctx context.Context, tx pgx.Tx, otpReq *model.OtpRequest) error {
	query := `
		UPDATE otp_requests
		SET attempts = $1, 
		    verified_at = $2, 
		    reset_token_hash = $3, 
		    token_expires_at = $4
		WHERE id = $5
	`
	_, err := tx.Exec(ctx, query,
		otpReq.Attempts,
		otpReq.VerifiedAt,
		otpReq.ResetTokenHash,
		otpReq.TokenExpiresAt,
		otpReq.ID,
	)
	return err
}

func (r *authRepository) GetOtpRequestByTokenHash(ctx context.Context, tx pgx.Tx, tokenHash string) (*model.OtpRequest, error) {
	query := `
		SELECT id, identifier, otp_hash, purpose, attempts, expires_at, verified_at, reset_token_hash, token_expires_at, created_at
		FROM otp_requests
		WHERE reset_token_hash = $1
		FOR UPDATE
	`
	var otpReq model.OtpRequest
	err := tx.QueryRow(ctx, query, tokenHash).Scan(
		&otpReq.ID,
		&otpReq.Identifier,
		&otpReq.OtpHash,
		&otpReq.Purpose,
		&otpReq.Attempts,
		&otpReq.ExpiresAt,
		&otpReq.VerifiedAt,
		&otpReq.ResetTokenHash,
		&otpReq.TokenExpiresAt,
		&otpReq.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &otpReq, nil
}

func (r *authRepository) InvalidateResetToken(ctx context.Context, tx pgx.Tx, id string) error {
	query := `
		UPDATE otp_requests
		SET reset_token_hash = NULL, token_expires_at = NULL
		WHERE id = $1
	`
	_, err := tx.Exec(ctx, query, id)
	return err
}

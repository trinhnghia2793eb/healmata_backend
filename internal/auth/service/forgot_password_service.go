package service

import (
	"context"
	"fmt"
	"log"
	"time"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	
	"healmata_backend/internal/auth/dto"
	authError "healmata_backend/internal/auth/errors"
	"healmata_backend/internal/auth/repository"
	"healmata_backend/internal/auth/model"
	"healmata_backend/internal/auth/otp"
	"healmata_backend/pkg/db"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *authService) ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequestDTO) (*dto.ForgotPasswordResponseDTO, error) {
	identifier := req.Identifier
	
	user, err := s.repo.GetUserByIdentifier(ctx, identifier)
	if err != nil || user == nil || user.ID == "" {
		return nil, authError.ErrUserNotFound
	}

	latestOtp, err := s.repo.GetLatestOtpRequest(ctx, identifier, "reset_password")
	if err != nil {
		return nil, authError.ErrForgotInternal
	}
	
	resendAfter := 60 * time.Second
	if latestOtp != nil {
		if time.Since(latestOtp.CreatedAt) < resendAfter {
			return nil, authError.ErrTooManyRequests
		}
	}

	otpCode, err := otp.GenerateOTP(6)
	if err != nil {
		return nil, authError.ErrForgotInternal
	}
	otpHash := otp.HashOTP(otpCode)

	expiresIn := 300 * time.Second
	payload := &repository.CreateOtpRequestPayload{
		Identifier: identifier,
		OtpHash:    otpHash,
		Purpose:    "reset_password",
		ExpiresAt:  time.Now().Add(expiresIn),
	}

	var newOtp *model.OtpRequest
	err = db.WithTransaction(ctx, s.dbPool, func(tx pgx.Tx) error {
		createdOtp, err := s.repo.CreateOtpRequest(ctx, tx, payload)
		if err != nil {
			return err
		}
		newOtp = createdOtp
		return nil
	})

	if err != nil {
		log.Printf("[ForgotPassword] Lỗi tạo OTP trong DB: %v", err)
		return nil, authError.ErrForgotInternal
	}

	go func(target, code string) {
		// TODO: 
		// log.Printf("\n========== OTP THÔNG BÁO ==========\n")
		// log.Printf("Gửi OTP [%s] tới [%s]", code, target)
		// log.Printf("===================================\n")

		// call SendOTP from pkg/email
		log.Printf("[Email Service] Đang gửi OTP tới %s...", target)
		err := s.emailSender.SendOTP(target, code)
		if err != nil {
			log.Printf("[Email Service ERROR] Không thể gửi OTP tới %s. Lỗi: %v", target, err)
			return
		}
		log.Printf("[Email Service] Đã gửi OTP thành công tới %s", target)
	}(identifier, otpCode)

	return &dto.ForgotPasswordResponseDTO{
		ResetRequestId: newOtp.ID,
		OtpLength:      6,
		ExpiresIn:      int64(expiresIn.Seconds()),
		ResendAfter:    int64(resendAfter.Seconds()),
	}, nil
}

func (s *authService) VerifyResetOtp(ctx context.Context, req *dto.VerifyResetOtpRequestDTO) (*dto.VerifyResetOtpResponseDTO, error) {
	var response *dto.VerifyResetOtpResponseDTO
	var businessErr error // business error flag (for return error without rollback transaction)

	err := db.WithTransaction(ctx, s.dbPool, func(tx pgx.Tx) error {
		otpReq, err := s.repo.GetOtpRequestByID(ctx, tx, req.ResetRequestId)
		if err != nil {
			return err
		}

		if otpReq == nil || otpReq.Purpose != "reset_password" {
			businessErr = authError.ErrInvalidOtp
			return nil
		}
		if otpReq.VerifiedAt != nil {
			businessErr = authError.ErrInvalidOtp
			return nil
		}
		if otpReq.Attempts >= 5 {
			businessErr = authError.ErrTooManyAttempts
			return nil
		}
		if time.Now().After(otpReq.ExpiresAt) {
			businessErr = authError.ErrExpiredOtp
			return nil
		}

		// get OTP from request
		inputOtpHash := otp.HashOTP(req.Otp) 
		// wrong OTP --> update attempt --> save into DB --> return nil
		if inputOtpHash != otpReq.OtpHash {
			otpReq.Attempts++
			if updateErr := s.repo.UpdateOtpRequest(ctx, tx, otpReq); updateErr != nil {
				return updateErr // Lỗi DB, Rollback
			}
			
			// assign businessErr + return nil --> transaction completed
			businessErr = authError.ErrInvalidOtp
			return nil 
		}
		// right OTP --> create reset token --> hash token --> update record
		tokenBytes := make([]byte, 32)
		if _, err := rand.Read(tokenBytes); err != nil {
			return err
		}
		resetTokenPlain := hex.EncodeToString(tokenBytes)

		tokenHash := sha256.Sum256([]byte(resetTokenPlain))
		resetTokenDBHash := hex.EncodeToString(tokenHash[:])

		now := time.Now()
		tokenExpiresIn := 600 * time.Second
		tokenExpiresAt := now.Add(tokenExpiresIn)

		otpReq.VerifiedAt = &now
		otpReq.ResetTokenHash = &resetTokenDBHash
		otpReq.TokenExpiresAt = &tokenExpiresAt

		if updateErr := s.repo.UpdateOtpRequest(ctx, tx, otpReq); updateErr != nil {
			return updateErr
		}

		response = &dto.VerifyResetOtpResponseDTO{
			ResetToken: resetTokenPlain,
			ExpiresIn:  int64(tokenExpiresIn.Seconds()),
		}

		return nil
	}) // end transaction

	// Handle internal error after transaction
	if err != nil {
		log.Printf("[VerifyResetOtp] Error: %v", err)
		return nil, authError.ErrOtpInternal
	}
	// Handle business error
	if businessErr != nil {
		return nil, businessErr
	}

	return response, nil
}

func (s *authService) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequestDTO) (*dto.ResetPasswordResponseDTO, error) {
	hash := sha256.Sum256([]byte(req.ResetToken))
	tokenHash := fmt.Sprintf("%x", hash[:])

	var resp dto.ResetPasswordResponseDTO

	err := db.WithTransaction(ctx, s.dbPool, func(tx pgx.Tx) error {	
		otpReq, err := s.repo.GetOtpRequestByTokenHash(ctx, tx, tokenHash)
		if err != nil {
			return authError.ErrResetPassInternal
		}

		if otpReq == nil || otpReq.ResetTokenHash == nil || *otpReq.ResetTokenHash == "" {
			return authError.ErrResetTokenExpired
		}
		if otpReq.TokenExpiresAt != nil && otpReq.TokenExpiresAt.Before(time.Now()) {
			return authError.ErrResetTokenExpired
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return authError.ErrResetPassInternal
		}
		err = s.repo.UpdateUserPassword(ctx, tx, otpReq.Identifier, string(hashedPassword))
		if err != nil {
			return authError.ErrResetPassInternal
		}
		err = s.repo.InvalidateResetToken(ctx, tx, otpReq.ID)
		if err != nil {
			return authError.ErrResetPassInternal
		}

		resp.PasswordReset = true
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &resp, nil
}
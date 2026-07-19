package errors

import (
	"net/http"
)

var (
	AUTH_REG_001 = NewAppError(http.StatusConflict, "AUTH_REG_001", "EMAIL_EXISTS")
	AUTH_REG_002 = NewAppError(http.StatusConflict, "AUTH_REG_002", "PHONE_EXISTS")
	AUTH_REG_003 = NewAppError(http.StatusUnprocessableEntity, "AUTH_REG_003", "INVALID_PASSWORD")
	AUTH_REG_004 = NewAppError(http.StatusUnprocessableEntity, "AUTH_REG_004", "INVALID_NAME")
	AUTH_REG_005 = NewAppError(http.StatusInternalServerError, "AUTH_REG_005", "INTERNAL_ERROR")
	AUTH_REG_006 = NewAppError(http.StatusUnprocessableEntity, "AUTH_REG_006", "INVALID_IDENTIFIER")
	AUTH_REG_007 = NewAppError(http.StatusUnprocessableEntity, "AUTH_REG_007", "INVALID_IDENTIFIER")
	AUTH_REG_008 = NewAppError(http.StatusUnprocessableEntity, "AUTH_REG_008", "CONFIRM_PASSWORD_MISMATCH")

	// Alias for compatibility/readability:
	ErrEmailExists             = AUTH_REG_001
	ErrPhoneExists             = AUTH_REG_002
	ErrInvalidPassword         = AUTH_REG_003
	ErrInvalidName             = AUTH_REG_004
	ErrInternalError           = AUTH_REG_005
	ErrInvalidEmail            = AUTH_REG_006
	ErrInvalidPhone            = AUTH_REG_007
	ErrConfirmPasswordMismatch = AUTH_REG_008

	// =======================================================================
	// FORGOT PASSWORD
	AUTH_FORGOT_001 = NewAppError(http.StatusNotFound, "AUTH_FORGOT_001", "USER_NOT_FOUND")
	AUTH_FORGOT_002 = NewAppError(http.StatusTooManyRequests, "AUTH_FORGOT_002", "TOO_MANY_REQUESTS")
	AUTH_FORGOT_003 = NewAppError(http.StatusInternalServerError, "AUTH_FORGOT_003", "INTERNAL_ERROR")

	ErrUserNotFound    = AUTH_FORGOT_001
	ErrTooManyRequests = AUTH_FORGOT_002
	ErrForgotInternal  = AUTH_FORGOT_003

	// =======================================================================
	// VERIFY RESET OTP
	AUTH_OTP_001 = NewAppError(http.StatusBadRequest, "AUTH_OTP_001", "INVALID_OTP")
	AUTH_OTP_002 = NewAppError(http.StatusBadRequest, "AUTH_OTP_002", "EXPIRED_OTP")
	AUTH_OTP_003 = NewAppError(http.StatusTooManyRequests, "AUTH_OTP_003", "TOO_MANY_ATTEMPTS")
	AUTH_OTP_004 = NewAppError(http.StatusInternalServerError, "AUTH_OTP_004", "INTERNAL_ERROR")

	ErrInvalidOtp      = AUTH_OTP_001
	ErrExpiredOtp      = AUTH_OTP_002
	ErrTooManyAttempts = AUTH_OTP_003
	ErrOtpInternal     = AUTH_OTP_004

	// =======================================================================
	// RESET PASSWORD
	AUTH_RESET_001 = NewAppError(http.StatusBadRequest, "AUTH_RESET_001", "PASSWORD_INVALID")
	AUTH_RESET_002 = NewAppError(http.StatusBadRequest, "AUTH_RESET_002", "PASSWORD_MISMATCH")
	AUTH_RESET_003 = NewAppError(http.StatusGone, "AUTH_RESET_003", "RESET_TOKEN_EXPIRED")
	AUTH_RESET_004 = NewAppError(http.StatusInternalServerError, "AUTH_RESET_004", "INTERNAL_ERROR")

	ErrPasswordInvalid   = AUTH_RESET_001
	ErrPasswordMismatch  = AUTH_RESET_002
	ErrResetTokenExpired = AUTH_RESET_003
	ErrResetPassInternal = AUTH_RESET_004

	// General Validation / JSON parsing error:
	ErrInvalidJSON = NewAppError(http.StatusBadRequest, "INVALID_JSON", "Định dạng JSON không hợp lệ")
)

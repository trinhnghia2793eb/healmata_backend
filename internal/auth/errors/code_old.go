package errors

// import (
// 	"net/http"
// )

// var (
// 	// =======================================================================
// 	// Shared Validation Errors (AUTH_VAL_*)
// 	AUTH_VAL_001 = NewAppError(http.StatusInternalServerError, "AUTH_VAL_001", "INTERNAL_ERROR")
// 	AUTH_VAL_002 = NewAppError(http.StatusBadRequest, "AUTH_VAL_002", "INVALID_JSON")
// 	AUTH_VAL_003 = NewAppError(http.StatusUnprocessableEntity, "AUTH_VAL_003", "INVALID_EMAIL")
// 	AUTH_VAL_004 = NewAppError(http.StatusUnprocessableEntity, "AUTH_VAL_004", "INVALID_PHONE")
// 	AUTH_VAL_005 = NewAppError(http.StatusUnprocessableEntity, "AUTH_VAL_005", "INVALID_PASSWORD")

// 	ErrInternalError   = AUTH_VAL_001
// 	ErrInvalidJSON     = AUTH_VAL_002
// 	ErrInvalidEmail    = AUTH_VAL_003
// 	ErrInvalidPhone    = AUTH_VAL_004
// 	ErrInvalidPassword = AUTH_VAL_005

// 	// =======================================================================
// 	// REGISTER (AUTH_REG_*)
// 	AUTH_REG_001 = NewAppError(http.StatusConflict, "AUTH_REG_001", "EMAIL_EXISTS")
// 	AUTH_REG_002 = NewAppError(http.StatusConflict, "AUTH_REG_002", "PHONE_EXISTS")
// 	AUTH_REG_003 = NewAppError(http.StatusUnprocessableEntity, "AUTH_REG_003", "INVALID_PASSWORD")
// 	AUTH_REG_004 = NewAppError(http.StatusUnprocessableEntity, "AUTH_REG_004", "INVALID_NAME")
// 	AUTH_REG_005 = NewAppError(http.StatusInternalServerError, "AUTH_REG_005", "INTERNAL_ERROR")
// 	AUTH_REG_006 = NewAppError(http.StatusServiceUnavailable, "AUTH_REG_006", "NETWORK_ERROR")
// 	AUTH_REG_007 = NewAppError(http.StatusUnprocessableEntity, "AUTH_REG_007", "PASSWORD_MISMATCH")

// 	ErrEmailExists             = AUTH_REG_001
// 	ErrPhoneExists             = AUTH_REG_002
// 	ErrInvalidPasswordReg      = AUTH_REG_003
// 	ErrInvalidName             = AUTH_REG_004
// 	ErrRegisterInternal        = AUTH_REG_005
// 	ErrConfirmPasswordMismatch = AUTH_REG_007

// 	// =======================================================================
// 	// LOGIN (AUTH_LOGIN_*)
// 	AUTH_LOGIN_001 = NewAppError(http.StatusUnauthorized, "AUTH_LOGIN_001", "INVALID_CREDENTIAL")
// 	AUTH_LOGIN_002 = NewAppError(http.StatusNotFound, "AUTH_LOGIN_002", "USER_NOT_FOUND")
// 	AUTH_LOGIN_003 = NewAppError(http.StatusForbidden, "AUTH_LOGIN_003", "USER_DISABLED")
// 	AUTH_LOGIN_004 = NewAppError(http.StatusTooManyRequests, "AUTH_LOGIN_004", "TOO_MANY_ATTEMPTS")
// 	AUTH_LOGIN_005 = NewAppError(http.StatusInternalServerError, "AUTH_LOGIN_005", "INTERNAL_ERROR")
// 	AUTH_LOGIN_006 = NewAppError(http.StatusServiceUnavailable, "AUTH_LOGIN_006", "NETWORK_ERROR")

// 	// =======================================================================
// 	// FORGOT PASSWORD
// 	AUTH_FORGOT_001 = NewAppError(http.StatusNotFound, "AUTH_FORGOT_001", "USER_NOT_FOUND")
// 	AUTH_FORGOT_002 = NewAppError(http.StatusTooManyRequests, "AUTH_FORGOT_002", "TOO_MANY_REQUESTS")
// 	AUTH_FORGOT_003 = NewAppError(http.StatusInternalServerError, "AUTH_FORGOT_003", "INTERNAL_ERROR")
// 	AUTH_FORGOT_004 = NewAppError(http.StatusServiceUnavailable, "AUTH_FORGOT_004", "NETWORK_ERROR")

// 	ErrUserNotFound    = AUTH_FORGOT_001
// 	ErrTooManyRequests = AUTH_FORGOT_002
// 	ErrForgotInternal  = AUTH_FORGOT_003

// 	// =======================================================================
// 	// VERIFY RESET OTP
// 	AUTH_OTP_001 = NewAppError(http.StatusBadRequest, "AUTH_OTP_001", "INVALID_OTP")
// 	AUTH_OTP_002 = NewAppError(http.StatusBadRequest, "AUTH_OTP_002", "EXPIRED_OTP")
// 	AUTH_OTP_003 = NewAppError(http.StatusTooManyRequests, "AUTH_OTP_003", "TOO_MANY_ATTEMPTS")
// 	AUTH_OTP_004 = NewAppError(http.StatusInternalServerError, "AUTH_OTP_004", "INTERNAL_ERROR")
// 	AUTH_OTP_005 = NewAppError(http.StatusServiceUnavailable, "AUTH_OTP_005", "NETWORK_ERROR")

// 	ErrInvalidOtp      = AUTH_OTP_001
// 	ErrExpiredOtp      = AUTH_OTP_002
// 	ErrTooManyAttempts = AUTH_OTP_003
// 	ErrOtpInternal     = AUTH_OTP_004

// 	// =======================================================================
// 	// RESET PASSWORD
// 	AUTH_RESET_001 = NewAppError(http.StatusBadRequest, "AUTH_RESET_001", "PASSWORD_INVALID")
// 	AUTH_RESET_002 = NewAppError(http.StatusBadRequest, "AUTH_RESET_002", "PASSWORD_MISMATCH")
// 	AUTH_RESET_003 = NewAppError(http.StatusGone, "AUTH_RESET_003", "RESET_TOKEN_EXPIRED")
// 	AUTH_RESET_004 = NewAppError(http.StatusInternalServerError, "AUTH_RESET_004", "INTERNAL_ERROR")
// 	AUTH_RESET_005 = NewAppError(http.StatusServiceUnavailable, "AUTH_RESET_005", "NETWORK_ERROR")

// 	ErrPasswordInvalid   = AUTH_RESET_001
// 	ErrPasswordMismatch  = AUTH_RESET_002
// 	ErrResetTokenExpired = AUTH_RESET_003
// 	ErrResetPassInternal = AUTH_RESET_004
// )

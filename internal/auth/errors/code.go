package errors

import (
	"net/http"
)

var (
	// Shared Validation Errors (AUTH_VAL_*)
	AUTH_VAL_001 = NewAppError(http.StatusBadRequest, "AUTH_VAL_001", "INVALID_JSON")
	AUTH_VAL_002 = NewAppError(http.StatusUnprocessableEntity, "AUTH_VAL_002", "INVALID_EMAIL")
	AUTH_VAL_003 = NewAppError(http.StatusUnprocessableEntity, "AUTH_VAL_003", "INVALID_PHONE")
	AUTH_VAL_004 = NewAppError(http.StatusUnprocessableEntity, "AUTH_VAL_004", "INVALID_PASSWORD")
	AUTH_VAL_005 = NewAppError(http.StatusInternalServerError, "AUTH_VAL_005", "INTERNAL_ERROR")

	// Register Errors (AUTH_REG_*)
	AUTH_REG_001 = NewAppError(http.StatusConflict, "AUTH_REG_001", "EMAIL_EXISTS")
	AUTH_REG_002 = NewAppError(http.StatusConflict, "AUTH_REG_002", "PHONE_EXISTS")
	AUTH_REG_003 = NewAppError(http.StatusUnprocessableEntity, "AUTH_REG_003", "INVALID_PASSWORD")
	AUTH_REG_004 = NewAppError(http.StatusUnprocessableEntity, "AUTH_REG_004", "INVALID_NAME")
	AUTH_REG_005 = NewAppError(http.StatusInternalServerError, "AUTH_REG_005", "INTERNAL_ERROR")
	AUTH_REG_006 = NewAppError(http.StatusServiceUnavailable, "AUTH_REG_006", "NETWORK_ERROR")
	AUTH_REG_007 = NewAppError(http.StatusUnprocessableEntity, "AUTH_REG_007", "PASSWORD_MISMATCH")

	// Login Errors (AUTH_LOGIN_*)
	AUTH_LOGIN_001 = NewAppError(http.StatusUnauthorized, "AUTH_LOGIN_001", "INVALID_CREDENTIAL")
	AUTH_LOGIN_002 = NewAppError(http.StatusNotFound, "AUTH_LOGIN_002", "USER_NOT_FOUND")
	AUTH_LOGIN_003 = NewAppError(http.StatusForbidden, "AUTH_LOGIN_003", "USER_DISABLED")
	AUTH_LOGIN_004 = NewAppError(http.StatusTooManyRequests, "AUTH_LOGIN_004", "TOO_MANY_ATTEMPTS")
	AUTH_LOGIN_005 = NewAppError(http.StatusInternalServerError, "AUTH_LOGIN_005", "INTERNAL_ERROR")
	AUTH_LOGIN_006 = NewAppError(http.StatusServiceUnavailable, "AUTH_LOGIN_006", "NETWORK_ERROR")

	// Aliases for compatibility/readability:
	ErrEmailExists             = AUTH_REG_001
	ErrPhoneExists             = AUTH_REG_002
	ErrInvalidPasswordReg      = AUTH_REG_003
	ErrInvalidName             = AUTH_REG_004
	ErrConfirmPasswordMismatch = AUTH_REG_007

	ErrInvalidEmail    = AUTH_VAL_002
	ErrInvalidPhone    = AUTH_VAL_003
	ErrInvalidPassword = AUTH_VAL_004
	ErrInternalError   = AUTH_VAL_005

	ErrInvalidJSON = AUTH_VAL_001
)

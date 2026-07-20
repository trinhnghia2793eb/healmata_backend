package errors

import (
	"net/http"
)

// =======================================================================
// SHARED VALIDATION ERROR
var Validation = struct {
	InternalError   *AppError
	InvalidJson     *AppError
	InvalidEmail    *AppError
	InvalidPhone    *AppError
	InvalidPassword *AppError
}{
	InternalError:   NewAppError(http.StatusInternalServerError, "AUTH_VAL_001", "INTERNAL_ERROR"),
	InvalidJson:     NewAppError(http.StatusBadRequest, "AUTH_VAL_002", "INVALID_JSON"),
	InvalidEmail:    NewAppError(http.StatusUnprocessableEntity, "AUTH_VAL_003", "INVALID_EMAIL"),
	InvalidPhone:    NewAppError(http.StatusUnprocessableEntity, "AUTH_VAL_004", "INVALID_PHONE"),
	InvalidPassword: NewAppError(http.StatusUnprocessableEntity, "AUTH_VAL_005", "INVALID_PASSWORD"),
}

// =======================================================================
// REGISTER (AUTH_REG_*)
var Register = struct {
	EmailExists             *AppError
	PhoneExists             *AppError
	InvalidPasswordReg      *AppError
	InvalidName             *AppError
	InternalError           *AppError
	NetworkError            *AppError
	ConfirmPasswordMismatch *AppError
}{
	EmailExists:             NewAppError(http.StatusConflict, "AUTH_REG_001", "EMAIL_EXISTS"),
	PhoneExists:             NewAppError(http.StatusConflict, "AUTH_REG_002", "PHONE_EXISTS"),
	InvalidPasswordReg:      NewAppError(http.StatusUnprocessableEntity, "AUTH_REG_003", "INVALID_PASSWORD"),
	InvalidName:             NewAppError(http.StatusUnprocessableEntity, "AUTH_REG_004", "INVALID_NAME"),
	InternalError:           NewAppError(http.StatusInternalServerError, "AUTH_REG_005", "INTERNAL_ERROR"),
	NetworkError:            NewAppError(http.StatusServiceUnavailable, "AUTH_REG_006", "NETWORK_ERROR"),
	ConfirmPasswordMismatch: NewAppError(http.StatusUnprocessableEntity, "AUTH_REG_007", "PASSWORD_MISMATCH"),
}

// =======================================================================
// LOGIN (AUTH_LOGIN_*)
var Login = struct {
	InvalidCredential *AppError
	UserNotFound      *AppError
	UserDisabled      *AppError
	TooManyAttempts   *AppError
	InternalError     *AppError
	NetworkError      *AppError
}{
	InvalidCredential: NewAppError(http.StatusUnauthorized, "AUTH_LOGIN_001", "INVALID_CREDENTIAL"),
	UserNotFound:      NewAppError(http.StatusNotFound, "AUTH_LOGIN_002", "USER_NOT_FOUND"),
	UserDisabled:      NewAppError(http.StatusForbidden, "AUTH_LOGIN_003", "USER_DISABLED"),
	TooManyAttempts:   NewAppError(http.StatusTooManyRequests, "AUTH_LOGIN_004", "TOO_MANY_ATTEMPTS"),
	InternalError:     NewAppError(http.StatusInternalServerError, "AUTH_LOGIN_005", "INTERNAL_ERROR"),
	NetworkError:      NewAppError(http.StatusServiceUnavailable, "AUTH_LOGIN_006", "NETWORK_ERROR"),
}

// =======================================================================
// FORGOT PASSWORD (AUTH_FORGOT_*)
var ForgotPassword = struct {
	UserNotFound    *AppError
	TooManyRequests *AppError
	InternalError   *AppError
	NetworkError    *AppError
}{
	UserNotFound:    NewAppError(http.StatusNotFound, "AUTH_FORGOT_001", "USER_NOT_FOUND"),
	TooManyRequests: NewAppError(http.StatusTooManyRequests, "AUTH_FORGOT_002", "TOO_MANY_REQUESTS"),
	InternalError:   NewAppError(http.StatusInternalServerError, "AUTH_FORGOT_003", "INTERNAL_ERROR"),
	NetworkError:    NewAppError(http.StatusServiceUnavailable, "AUTH_FORGOT_004", "NETWORK_ERROR"),
}

// =======================================================================
// VERIFY OTP (AUTH_OTP_*)
var VerifyOtp = struct {
	InvalidOtp      *AppError
	ExpiredOtp      *AppError
	TooManyAttempts *AppError
	InternalError   *AppError
	NetworkError    *AppError
}{
	InvalidOtp:      NewAppError(http.StatusBadRequest, "AUTH_OTP_001", "INVALID_OTP"),
	ExpiredOtp:      NewAppError(http.StatusBadRequest, "AUTH_OTP_002", "EXPIRED_OTP"),
	TooManyAttempts: NewAppError(http.StatusTooManyRequests, "AUTH_OTP_003", "TOO_MANY_ATTEMPTS"),
	InternalError:   NewAppError(http.StatusInternalServerError, "AUTH_OTP_004", "INTERNAL_ERROR"),
	NetworkError:    NewAppError(http.StatusServiceUnavailable, "AUTH_OTP_005", "NETWORK_ERROR"),
}

// =======================================================================
// RESET PASSWORD (AUTH_RESET_*)
var ResetPassword = struct {
	PasswordInvalid   *AppError
	PasswordMismatch  *AppError
	ResetTokenExpired *AppError
	InternalError     *AppError
	NetworkError      *AppError
}{
	PasswordInvalid:   NewAppError(http.StatusBadRequest, "AUTH_RESET_001", "PASSWORD_INVALID"),
	PasswordMismatch:  NewAppError(http.StatusBadRequest, "AUTH_RESET_002", "PASSWORD_MISMATCH"),
	ResetTokenExpired: NewAppError(http.StatusGone, "AUTH_RESET_003", "RESET_TOKEN_EXPIRED"),
	InternalError:     NewAppError(http.StatusInternalServerError, "AUTH_RESET_004", "INTERNAL_ERROR"),
	NetworkError:      NewAppError(http.StatusServiceUnavailable, "AUTH_RESET_005", "NETWORK_ERROR"),
}

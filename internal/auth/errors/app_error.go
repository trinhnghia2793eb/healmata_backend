package errors

import (
	"net/http"

	"healmata_backend/pkg/response"
)

// =======================================================================
// SHARED VALIDATION ERROR
var Validation = struct {
	InternalError *response.AppError
	InvalidJson   *response.AppError
}{
	InternalError: response.NewAppError(http.StatusInternalServerError, "AUTH_VAL_001", "INTERNAL_ERROR"),
	InvalidJson:   response.NewAppError(http.StatusBadRequest, "AUTH_VAL_002", "INVALID_JSON"),
}

// =======================================================================
// REGISTER (AUTH_REG_*)
var Register = struct {
	EmailExists *response.AppError
	PhoneExists *response.AppError
	// InvalidPasswordReg      *response.AppError
	// InvalidName             *response.AppError
	InternalError *response.AppError
	NetworkError  *response.AppError
	// ConfirmPasswordMismatch *response.AppError
}{
	EmailExists: response.NewAppError(http.StatusConflict, "AUTH_REG_001", "EMAIL_EXISTS"),
	PhoneExists: response.NewAppError(http.StatusConflict, "AUTH_REG_002", "PHONE_EXISTS"),
	// InvalidPasswordReg:      response.NewAppError(http.StatusUnprocessableEntity, "AUTH_REG_003", "INVALID_PASSWORD"),
	// InvalidName:             response.NewAppError(http.StatusUnprocessableEntity, "AUTH_REG_004", "INVALID_NAME"),
	InternalError: response.NewAppError(http.StatusInternalServerError, "AUTH_REG_005", "INTERNAL_ERROR"),
	NetworkError:  response.NewAppError(http.StatusServiceUnavailable, "AUTH_REG_006", "NETWORK_ERROR"),
	// ConfirmPasswordMismatch: response.NewAppError(http.StatusUnprocessableEntity, "AUTH_REG_007", "PASSWORD_MISMATCH"),
}

// =======================================================================
// LOGIN (AUTH_LOGIN_*)
var Login = struct {
	InvalidCredential *response.AppError
	UserNotFound      *response.AppError
	UserDisabled      *response.AppError
	TooManyAttempts   *response.AppError
	InternalError     *response.AppError
	NetworkError      *response.AppError
}{
	InvalidCredential: response.NewAppError(http.StatusUnauthorized, "AUTH_LOGIN_001", "INVALID_CREDENTIAL"),
	UserNotFound:      response.NewAppError(http.StatusNotFound, "AUTH_LOGIN_002", "USER_NOT_FOUND"),
	UserDisabled:      response.NewAppError(http.StatusForbidden, "AUTH_LOGIN_003", "USER_DISABLED"),
	TooManyAttempts:   response.NewAppError(http.StatusTooManyRequests, "AUTH_LOGIN_004", "TOO_MANY_ATTEMPTS"),
	InternalError:     response.NewAppError(http.StatusInternalServerError, "AUTH_LOGIN_005", "INTERNAL_ERROR"),
	NetworkError:      response.NewAppError(http.StatusServiceUnavailable, "AUTH_LOGIN_006", "NETWORK_ERROR"),
}

// =======================================================================
// FORGOT PASSWORD (AUTH_FORGOT_*)
var ForgotPassword = struct {
	UserNotFound      *response.AppError
	TooManyRequests   *response.AppError
	PhoneNotSupported *response.AppError
	InternalError     *response.AppError
	NetworkError      *response.AppError
}{
	UserNotFound:      response.NewAppError(http.StatusNotFound, "AUTH_FORGOT_001", "USER_NOT_FOUND"),
	TooManyRequests:   response.NewAppError(http.StatusTooManyRequests, "AUTH_FORGOT_002", "TOO_MANY_REQUESTS"),
	PhoneNotSupported: response.NewAppError(http.StatusBadRequest, "AUTH_FORGOT_PNS", "PHONE_NOT_SUPPORTED"),
	InternalError:     response.NewAppError(http.StatusInternalServerError, "AUTH_FORGOT_003", "INTERNAL_ERROR"),
	NetworkError:      response.NewAppError(http.StatusServiceUnavailable, "AUTH_FORGOT_004", "NETWORK_ERROR"),
}

// =======================================================================
// VERIFY OTP (AUTH_OTP_*)
var VerifyOtp = struct {
	InvalidOtp      *response.AppError
	ExpiredOtp      *response.AppError
	TooManyAttempts *response.AppError
	InternalError   *response.AppError
	NetworkError    *response.AppError
}{
	InvalidOtp:      response.NewAppError(http.StatusBadRequest, "AUTH_OTP_001", "INVALID_OTP"),
	ExpiredOtp:      response.NewAppError(http.StatusBadRequest, "AUTH_OTP_002", "EXPIRED_OTP"),
	TooManyAttempts: response.NewAppError(http.StatusTooManyRequests, "AUTH_OTP_003", "TOO_MANY_ATTEMPTS"),
	InternalError:   response.NewAppError(http.StatusInternalServerError, "AUTH_OTP_004", "INTERNAL_ERROR"),
	NetworkError:    response.NewAppError(http.StatusServiceUnavailable, "AUTH_OTP_005", "NETWORK_ERROR"),
}

// =======================================================================
// RESET PASSWORD (AUTH_RESET_*)
var ResetPassword = struct {
	// PasswordInvalid   *response.AppError
	// PasswordMismatch  *response.AppError
	ResetTokenExpired *response.AppError
	InternalError     *response.AppError
	NetworkError      *response.AppError
}{
	// PasswordInvalid:   response.NewAppError(http.StatusBadRequest, "AUTH_RESET_001", "PASSWORD_INVALID"),
	// PasswordMismatch:  response.NewAppError(http.StatusBadRequest, "AUTH_RESET_002", "PASSWORD_MISMATCH"),
	ResetTokenExpired: response.NewAppError(http.StatusGone, "AUTH_RESET_003", "RESET_TOKEN_EXPIRED"),
	InternalError:     response.NewAppError(http.StatusInternalServerError, "AUTH_RESET_004", "INTERNAL_ERROR"),
	NetworkError:      response.NewAppError(http.StatusServiceUnavailable, "AUTH_RESET_005", "NETWORK_ERROR"),
}

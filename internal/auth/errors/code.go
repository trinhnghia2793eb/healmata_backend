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

	// General Validation / JSON parsing error:
	ErrInvalidJSON = NewAppError(http.StatusBadRequest, "INVALID_JSON", "Định dạng JSON không hợp lệ")
)

package errors

import "errors"

const (
	AuthInternalError   = "AUTH_INTERNAL_ERROR"
	AuthBadRequest      = "AUTH_BAD_REQUEST"
	AuthUnauthorized    = "AUTH_UNAUTHORIZED"
	AuthForbidden       = "AUTH_FORBIDDEN"
	AuthValidationError = "AUTH_VALIDATION_ERROR"
)

var (
	ErrEmailExists   = errors.New("EMAIL_EXISTS")
	ErrPhoneExists   = errors.New("PHONE_EXISTS")
	ErrInternalError = errors.New("INTERNAL_ERROR")
)

type AppError struct {
	HTTPStatus int    `json:"-"`
	ErrorCode  string `json:"code"`
	Message    string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(httpStatus int, code string, message string) *AppError {
	return &AppError{
		HTTPStatus: httpStatus,
		ErrorCode:  code,
		Message:    message,
	}
}

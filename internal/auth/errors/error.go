package errors

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

func GetValidationErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Trường này không được bỏ trống"
	case "min":
		return "Độ dài ngắn hơn quy định tối thiểu (yêu cầu tối thiểu " + fe.Param() + " ký tự)"
	case "max":
		return "Độ dài vượt quá quy định tối đa (tối đa " + fe.Param() + " ký tự)"
	case "eqfield":
		return "Mật khẩu xác nhận không khớp"
	case "email":
		return "Định dạng email không hợp lệ"
	}
	return "Dữ liệu không hợp lệ"
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
type AppError struct {
	HTTPStatus int           `json:"-"`
	ErrorCode  string        `json:"code"`
	Message    string        `json:"message"`
	Details    []ErrorDetail `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(httpStatus int, errorCode string, message string) *AppError {
	return &AppError{
		HTTPStatus: httpStatus,
		ErrorCode:  errorCode,
		Message:    message,
	}
}

func NewValidationError(details []ErrorDetail) *AppError {
	return &AppError{
		HTTPStatus: http.StatusBadRequest,
		ErrorCode:  "VALIDATION_ERROR",
		Message:    "Dữ liệu đầu vào không hợp lệ",
		Details:    details,
	}
}

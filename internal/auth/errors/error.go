package errors

import (
	// "net/http"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
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

type ErrorDetail struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
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

// return error in middleware, handler, service?
func ReturnAppError(c *gin.Context, appErr *AppError) {
	c.AbortWithStatusJSON(appErr.HTTPStatus, gin.H{
		"success": false,
		"error":   appErr,
	})
}

func NewValidationError(details []ErrorDetail) *AppError {
	return &AppError{
		HTTPStatus: http.StatusBadRequest,
		ErrorCode:  "VALIDATION_ERROR",
		Message:    "Dữ liệu đầu vào không hợp lệ",
		Details:    details,
	}
}

// ParseValidationErrors parses erorrs from validator/v10
func ReturnValidationError(c *gin.Context, err error) {
	var details []ErrorDetail

	// convert erorr type into validator.ValidationErrors
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			// e.Field() --> field name (e.g: Email)
			// e.Tag() --> rule validation tag in DTO (e.g: required)
			rule := GetValidationRule(e.Field(), e.Tag())

			details = append(details, ErrorDetail{
				Field:   e.Field(),
				Code:    rule.Code,
				Message: rule.Message,
			})
		}
	} else {
		details = append(details, ErrorDetail{
			Field:   "payload",
			Code:    "VR-ERR",
			Message: "Cấu trúc dữ liệu đầu vào không hợp lệ",
		})
	}

	if len(details) > 0 {
		appErr := NewValidationError(details)
		ReturnAppError(c, appErr)
		return
	}

	ReturnAppError(c, Validation.InvalidJson)
}

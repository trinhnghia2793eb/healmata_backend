package errors

import (
	"fmt"
	"healmata_backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ValidationRule struct {
	Code    string
	Message string
}

var validationErrorMap = map[string]ValidationRule{
	"fullName.required":    {"VR-NAME-001", "Vui lòng nhập họ và tên"},
	"fullName.min":         {"VR-NAME-002", "Họ và tên quá ngắn"},
	"fullName.max":         {"VR-NAME-003", "Họ và tên quá dài"},
	"fullName.is_fullname": {"VR-NAME-004", "Họ và tên không hợp lệ"},

	"identifier.required":      {"VR-ID-001", "Vui lòng nhập email / sdt"},
	"identifier.is_identifier": {"VR-ID-002", "Thông tin đăng nhập không hợp lệ"},
	"identifier.max":           {"VR-ID-003", "Dữ liệu quá dài"},

	// "register.email.required": {"VR-REG-001", "Email không được để trống"},
	// "register.email.email":    {"VR-REG-002", "Email không đúng định dạng"},
	// "register.name.required":  {"VR-REG-003", "Họ và tên không được để trống"},
	// "register.phone.required": {"VR-REG-006", "Số điện thoại không được để trống"},

	"password.required": {"VR-PWD-001", "Vui lòng nhập mật khẩu"},
	"password.min":      {"VR-PWD-002", "Mật khẩu phải có ít nhất 8 ký tự"},
	"password.max":      {"VR-PWD-003", "Mật khẩu quá dài"},

	"confirmPassword.required": {"VR-CONFIRM-001", "Vui lòng xác nhận mật khẩu"},
	"confirmPassword.eqfield":  {"VR-CONFIRM-002", "Mật khẩu xác nhận không khớp"},

	"resetRequestId.required": {"VR-RESET-REQUEST-ID-001", "Thiếu mã yêu cầu đặt lại mật khẩu"},
	"resetRequestId.uuid":     {"VR-RESET-REQUEST-ID-002", "Mã yêu cầu đặt lại mật khẩu không hợp lệ"},

	"otp.required": {"VR-OTP-001", "Vui lòng nhập mã xác nhận"},
	"otp.numeric":  {"VR-OTP-002", "Mã xác nhận không hợp lệ"},
	"otp.len":      {"VR-OTP-003", "Mã xác nhận phải gồm 6 số"},

	"resetToken.required": {"VR-RESET-TOKEN-001", "Thiếu mã đặt lại mật khẩu"},

	"newPassword.required": {"VR-NEW-PWD-001", "Vui lòng nhập mật khẩu"},
	"newPassword.min":      {"VR-NEW-PWD-002", "Mật khẩu phải có ít nhất 8 ký tự"},
	"newPassword.max":      {"VR-NEW-PWD-003", "Mật khẩu quá dài"},
}

// get validation rule from validationErrorMap
func GetValidationRule(field, tag string) ValidationRule {
	key := fmt.Sprintf("%s.%s", field, tag)
	if rule, exists := validationErrorMap[key]; exists {
		return rule
	}
	return ValidationRule{
		Code:    "VR-ERR",
		Message: fmt.Sprintf("%s không hợp lệ (Lỗi: %s)", field, tag),
	}
}

// handle error & call pkg/response
func HandleValidationError(c *gin.Context, err error) {
	var details []response.ErrorDetail

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			rule := GetValidationRule(e.Field(), e.Tag())

			// push detail into ErrorDetail array
			details = append(details, response.ErrorDetail{
				Field:   e.Field(),
				Code:    rule.Code,
				Message: rule.Message,
			})
		}
	} else {
		// fallback
		details = append(details, response.ErrorDetail{
			Field:   "payload",
			Code:    "VR-INVALID-JSON",
			Message: "Cấu trúc dữ liệu đầu vào không hợp lệ",
		})
	}

	// create AppError struct
	appErr := &response.AppError{
		HTTPStatus: http.StatusBadRequest,
		ErrorCode:  "VALIDATION_ERROR",
		Message:    "Dữ liệu đầu vào không hợp lệ",
		Details:    details,
	}

	// call response->AbortWithError
	response.AbortWithError(c, appErr)
}

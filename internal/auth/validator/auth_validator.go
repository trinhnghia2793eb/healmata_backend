package validator

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	// Regex: Chỉ cho phép chữ cái (gồm tiếng Việt có dấu) và khoảng trắng. Không số, không ký tự đặc biệt.
	fullNameRegex = regexp.MustCompile(`^[a-zA-ZÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚĂĐĨŨƠàáâãèéêìíòóôõùúăđĩũơƯĂẠẢẤẦẨẪẬẮẰẲẴẶẸẺẼỀỀỂưăạảấầẩẫậắằẳẵặẹẻẽềềểỄỆỈỊỌỎỐỒỔỖỘỚỜỞỠỢỤỦỨỪễệỉịọỏốồổỗộớờởỡợụủứừỬỮỰỲỴÝỶỸửữựỳỵỷỹ\s]+$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex    = regexp.MustCompile(`^\+?[0-9]{9,15}$`)
)

// ValidateFullName chặn chuỗi rỗng sau khi trim và chặn ký tự đặc biệt
func ValidateFullName(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	trimmed := strings.TrimSpace(val)
	if len(trimmed) < 2 {
		return false
	}
	return fullNameRegex.MatchString(trimmed)
}

// ValidateIdentifier check xem nó là email hay phone hợp lệ không
func ValidateIdentifier(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	trimmed := strings.TrimSpace(val)

	if strings.Contains(trimmed, "@") {
		return emailRegex.MatchString(trimmed)
	}
	return phoneRegex.MatchString(trimmed)
}

// RegisterCustomValidators được gọi lúc khởi tạo server (ví dụ trong main.go)
func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("is_fullname", ValidateFullName)
		_ = v.RegisterValidation("is_identifier", ValidateIdentifier)
	}
}

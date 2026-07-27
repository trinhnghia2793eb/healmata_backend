package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ======================================================================= 
// App error
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

// ======================================================================= 
// SuccessReponse & Error Response
type SuccessResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Success bool      `json:"success"`
	Error   *AppError `json:"error"`
}

// return success - used in handler
func Success(c *gin.Context, data any, message string) {
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// return error - used in handler
func Error(c *gin.Context, appErr *AppError) {
	c.JSON(appErr.HTTPStatus, ErrorResponse{
		Success: false,
		Error:   appErr,
	})
}

// return error & stop middleware - used in middleware
func AbortWithError(c *gin.Context, appErr *AppError) {
	c.AbortWithStatusJSON(appErr.HTTPStatus, ErrorResponse{
		Success: false,
		Error:   appErr,
	})
}

package response

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type ErrorDetail struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

type ErrorResponse struct {
	Success bool `json:"success"`
	Error ErrorDetail `json:"error"`
}

func NewSuccessResponse(data any, message string) SuccessResponse {
    return SuccessResponse{
        Success: true,
		Data: data,
		Message: message,
    }
}

func NewErrorResponse(code string, message string) ErrorResponse {
    return ErrorResponse{
        Success: false,
		Error: ErrorDetail{
			Code: code,
			Message: message,
		},
    }
}
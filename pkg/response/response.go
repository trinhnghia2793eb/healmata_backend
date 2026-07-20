package response

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type ErrorResponse struct {
	Success bool `json:"success"`

	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

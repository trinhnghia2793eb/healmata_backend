package dto

type ForgotPasswordRequestDTO struct {
	Identifier string `json:"identifier" binding:"required,is_identifier"`
}

type ForgotPasswordResponseDTO struct {
	ResetRequestId string `json:"resetRequestId"`
	OtpLength      int    `json:"otpLength"`
	ExpiresIn      int64  `json:"expiresIn"`
	ResendAfter    int64  `json:"resendAfter"`
}

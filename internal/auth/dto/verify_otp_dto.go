package dto

type VerifyResetOtpRequestDTO struct {
	ResetRequestId string `json:"resetRequestId" binding:"required,uuid"`
	Otp string `json:"otp" binding:"required,numeric,len=6"`
}

type VerifyResetOtpResponseDTO struct {
	ResetToken string `json:"resetToken"`
	ExpiresIn  int64  `json:"expiresIn"`
}
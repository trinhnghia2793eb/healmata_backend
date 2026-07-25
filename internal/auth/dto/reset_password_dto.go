package dto

type ResetPasswordRequestDTO struct {
	ResetToken      string `json:"resetToken" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8,max=128"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=NewPassword"`
}

type ResetPasswordResponseDTO struct {
	PasswordReset bool `json:"passwordReset"`
}

package dto

type Device struct {
	Platform string `json:"platform"`
	DeviceID string `json:"deviceId"`
}

type SocialLoginRequestDTO struct {
	Provider      string  `json:"provider" binding:"required,oneof=google apple"`
	ProviderToken string  `json:"providerToken" binding:"required"`
	Device        *Device `json:"device"`
}

type SocialLoginResponseDTO struct {
	TokenResponseDTO
	IsNewUser     bool `json:"isNewUser"`
	LinkedAccount bool `json:"linkedAccount"`
}

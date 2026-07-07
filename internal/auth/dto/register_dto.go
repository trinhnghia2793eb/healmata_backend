package dto

type RegisterRequestDTO struct {
	FullName        string `json:"fullName" validate:"required, min=2,max=100"`
	Identifier      string `json:"identifier" validate:"required"`
	Password        string `json:"password" validate:"required, min=8, max=128"`
	ConfirmPassword string `json:"confirmPassword" validate:"required, eqfield=Password"`
}

type RegisterResponseDTO struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}
